/*
	VNFM SDK for Open Baton Managers
 */
package vnfmsdk

import (
	"os"
	"github.com/BurntSushi/toml"
	"encoding/json"
	"github.com/openbaton/go-openbaton/sdk"
	"github.com/openbaton/go-openbaton/catalogue"
	"os/signal"
)

// The VNFM config struct
type VnfmConfig struct {
	Type        string `toml:"type"`
	Endpoint    string `toml:"endpoint"`
	Description string `toml:"description"`
	Workers     int    `toml:"workers"`
	Username    string `toml:"username"`
	Password    string `toml:"password"`
	Allocate    bool   `toml:"allocate"`
	LogLevel    string `toml:"logLevel"`
	BrokerIp    string `toml:"brokerIp"`
	BrokerPort  int    `toml:"brokerPort"`
}

// Start the VNFM with config file
func Start(confPath string, h HandlerVnfm, name string) (error) {
	cfg := VnfmConfig{
		Type:        "unknown",
		Workers:     5,
		Allocate:    false,
		Description: "The Vnfm written in go",
		Username:    "openbaton-manager-user",
		Password:    "openbaton",
		LogLevel:    "DEBUG",
		BrokerIp:    "localhost",
		BrokerPort:  5672,
	}
	cfg.Endpoint = cfg.Type

	reader, err := os.Open(confPath)
	defer reader.Close()
	if err != nil {
		return err
	}
	if _, err := toml.DecodeReader(reader, &cfg); err != nil {
		return err
	}

	return startWithCfg(cfg, name, h)
}

// Start the VNFM with specific config
func StartWithConfig(typ, description, username, password, loglevel, brokerIp string, brokerPort, workers int, allocate bool, h HandlerVnfm, name string) (error) {
	cfg := VnfmConfig{
		Type:        typ,
		Workers:     workers,
		Allocate:    allocate,
		Description: description,
		Username:    username,
		Password:    password,
		LogLevel:    loglevel,
		BrokerIp:    brokerIp,
		BrokerPort:  brokerPort,
	}
	cfg.Endpoint = cfg.Type

	return startWithCfg(cfg, name, h)
}

func startWithCfg(cfg VnfmConfig, name string, h HandlerVnfm) error {
	logger := sdk.GetLogger(cfg.Type, cfg.LogLevel)
	logger.Infof("Starting VNFM of type %s", cfg.Type)
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	logger.Debugf("Config are %s", jsonCfg)

	endpoint := catalogue.Endpoint{
		Type:         cfg.Type,
		Endpoint:     cfg.Endpoint,
		Active:       true,
		Description:  cfg.Description,
		Enabled:      true,
		EndpointType: "RABBIT",
	}
	rabbitCredentials, err := sdk.GetVnfmCreds(cfg.Username, cfg.Password, cfg.BrokerIp, cfg.BrokerPort, &endpoint, "DEBUG")

	if err != nil {
		logger.Errorf("Error getting credentials: %v", err)
		return err
	}

	manager, err := sdk.NewManager(
		h,
		rabbitCredentials.RabbitUsername,
		rabbitCredentials.RabbitPassword,
		cfg.BrokerIp,
		cfg.BrokerPort,
		"openbaton-exchange",
		endpoint.Endpoint,
		cfg.Workers,
		cfg.Allocate,
		name,
		handleNfvMessage,
		"DEBUG",
		nil,
		nil,
	)
	if err != nil {
		logger.Errorf("Error while creating vnfm: %v", err)
		return err
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			logger.Infof("Received ctrl-c, unregistering")
			manager.Unregister(cfg.Type, rabbitCredentials.RabbitUsername, rabbitCredentials.RabbitPassword, &endpoint)
			go manager.Shutdown()
			os.Exit(0)
		}
	}()

	manager.Serve()

	return err
}
