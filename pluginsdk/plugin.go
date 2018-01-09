/*
	Plugin SDK for Open Baton Managers
 */
package pluginsdk

import (
	"os"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/openbaton/go-openbaton/sdk"
	"os/signal"
	"encoding/json"
)

// The Config struct for a plugin
type PluginConfig struct {
	Type       string `toml:"type"`
	Workers    int    `toml:"workers"`
	Username   string `toml:"username"`
	Password   string `toml:"password"`
	LogLevel   string `toml:"logLevel"`
	BrokerIp   string `toml:"brokerIp"`
	BrokerPort int    `toml:"brokerPort"`
}

// Start the plugin using the configuration file
func Start(confPath string, h HandlerVim, name string) (error) {
	cfg := PluginConfig{
		Type:       "unknown",
		Workers:    5,
		Username:   "openbaton-manager-user",
		Password:   "openbaton",
		LogLevel:   "DEBUG",
		BrokerIp:   "localhost",
		BrokerPort: 5672,
	}
	reader, err := os.Open(confPath)
	defer reader.Close()
	if err != nil {
		return err
	}
	if _, err := toml.DecodeReader(reader, &cfg); err != nil {
		return err
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: while loading config file %s: %v\n", confPath, err)
		os.Exit(100)
	}

	return startWithCfg(cfg, h, name)
}

// Start the plugin with specific configuration
func StartWithConfig(typ, username, password, loglevel, brokerip string, workers, brokerPort int, h HandlerVim, name string) (error) {
	cfg := PluginConfig{
		Type:       typ,
		Workers:    workers,
		Username:   username,
		Password:   password,
		LogLevel:   loglevel,
		BrokerIp:   brokerip,
		BrokerPort: brokerPort,
	}

	return startWithCfg(cfg, h, name)
}

func startWithCfg(cfg PluginConfig, h HandlerVim, name string) error {
	pluginId := fmt.Sprintf("vim-drivers.%s.%s", cfg.Type, name)
	logger := sdk.GetLogger(cfg.Type, cfg.LogLevel)
	logger.Infof("Starting Plugin of type %s", cfg.Type)
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	logger.Debugf("Config are %s", jsonCfg)
	rabbitCredentials, err := sdk.GetPluginCreds(cfg.Username, cfg.Password, cfg.BrokerIp, cfg.BrokerPort, pluginId, "DEBUG")

	if err != nil {
		logger.Errorf("Error getting credentials: %v", err)
		return err
	}

	manager, err := sdk.NewPluginManager(
		rabbitCredentials.RabbitUsername,
		rabbitCredentials.RabbitPassword,
		cfg.BrokerIp,
		cfg.BrokerPort,
		"openbaton-exchange",
		pluginId,
		cfg.Workers,
		handlePluginRequest,
		"DEBUG",
	)
	if err != nil {
		return err
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logger.Infof("Received ctrl-c, unregistering")
			manager.Unregister(cfg.Type, rabbitCredentials.RabbitUsername, rabbitCredentials.RabbitPassword)
			go manager.Shutdown()
			logger.Infof("Done")
			os.Exit(0)
		}
	}()

	wk := &worker{
		l: logger,
		h: h,
	}
	manager.Serve(wk)

	return err
}
