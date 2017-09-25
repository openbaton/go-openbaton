package vnfmsdk

import (
	"os"
	"github.com/BurntSushi/toml"
	"github.com/openbaton/go-openbaton/sdk"
	"github.com/openbaton/go-openbaton/catalogue"
	"os/signal"
)

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

func Start(confPath string, h HandlerVnfm, name string) (error) {
	cfg := VnfmConfig{
		Type:        "unknown",
		Description: "The Vnfm written in go",
		Workers:     5,
		Allocate:    false,
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

	logger := sdk.GetLogger(cfg.Type, cfg.LogLevel)
	logger.Infof("Starting VNFM of type %s", cfg.Type)

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

	manager, err := sdk.NewVnfmManager(
		rabbitCredentials.RabbitUsername,
		rabbitCredentials.RabbitPassword,
		cfg.BrokerIp,
		cfg.BrokerPort,
		"openbaton-exchange",
		endpoint.Endpoint,
		cfg.Workers,
		name,
		handleNfvMessage,
		"DEBUG",
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
			manager.Manager.Unregister(cfg.Type, rabbitCredentials.RabbitUsername, rabbitCredentials.RabbitPassword)
			go manager.Shutdown()
			os.Exit(0)
		}
	}()

	wk := &worker{
		l:        logger,
		Channel:  manager.Channel,
		handler:  h,
		Allocate: cfg.Allocate,
	}
	manager.Serve(wk)

	return err
}
