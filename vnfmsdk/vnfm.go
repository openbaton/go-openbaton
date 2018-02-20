//VNFM SDK for Open Baton VNFManagers. Uses the sdk package passing specific implementation of certain functions.
package vnfmsdk

import (
	"encoding/json"
	"os"
	"os/signal"

	"github.com/BurntSushi/toml"
	"github.com/openbaton/go-openbaton/catalogue"
	"github.com/openbaton/go-openbaton/sdk"
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
	Timeout     int    `toml:"timeout"`
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
		Timeout:     2,
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
func StartWithConfig(typ, description, username, password, loglevel, brokerIp string, brokerPort, workers, timeout int, allocate bool, h HandlerVnfm, name string) (error) {
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
		Timeout:     timeout,
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
	rabbitCredentials, err := sdk.GetVnfmCreds(cfg.Username, cfg.Password, cfg.BrokerIp, cfg.BrokerPort, cfg.Timeout, &endpoint, "DEBUG")

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
		sdk.OpenbatonExchangeName,
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
		for range c {
			logger.Infof("Received ctrl-c, unregistering")
			manager.Unregister(cfg.Type, rabbitCredentials.RabbitUsername, rabbitCredentials.RabbitPassword, &endpoint)
			go manager.Shutdown()
			os.Exit(0)
		}
	}()

	manager.Serve()

	return err
}
