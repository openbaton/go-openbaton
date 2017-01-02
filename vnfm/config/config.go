package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

// Config represents a generic config type for a VNFM,
// exporting some basic variables from the '[vnfm]' section of
// the config file.
type Config struct {
	Allocate bool

	LogColors bool
	LogFile   string
	LogLevel  log.Level

	Type        string
	Endpoint    string
	Description string

	// Properties contain the raw Properties from which this config has
	// been extracted. They also may contain implementation specific fields that
	// may be needed.
	Properties Properties
}

// Load loads a Config from an io.Reader containing TOML data.
func Load(reader io.Reader) (*Config, error) {
	props := make(Properties)

	if _, err := toml.DecodeReader(reader, &props); err != nil {
		return nil, err
	}

	return New(props)
}

func LoadFile(fileName string) (*Config, error) {
	reader, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return Load(reader)
}

func New(props Properties) (*Config, error) {
	vnfm, ok := props.Section("vnfm")
	if !ok {
		return nil, errors.New("malformed config - missing '[vnfm]' section")
	}

	allocate, _ := vnfm.ValueBool("allocate", true)

	vnfmType, set := vnfm.ValueString("type", "")
	if !set {
		return nil, errors.New("no vnfm.type in config")
	}

	endpoint, set := vnfm.ValueString("endpoint", "")
	if !set {
		return nil, errors.New("no vnfm.endpoint in config")
	}

	descr, _ := vnfm.ValueString("description", "")

	logColors := true
	logFile := ""
	logLevel := "WARN"

	logger, ok := vnfm.Section("logger")
	if ok {
		logColors, _ = logger.ValueBool("use-colors", logColors)
		logFile, _ = logger.ValueString("out-file", logFile)
		logLevel, _ = logger.ValueString("level", logLevel)
	}

	lvl, err := toLogLevel(logLevel)
	if err != nil {
		return nil, err
	}

	return &Config{
		Allocate:    allocate,
		Description: descr,
		Endpoint:    endpoint,
		LogColors:   logColors,
		LogFile:     logFile,
		LogLevel:    lvl,
		Properties:  props,
		Type:        vnfmType,
	}, nil
}

func toLogLevel(lvlStr string) (lvl log.Level, err error) {
	switch strings.ToUpper(lvlStr) {
	case "DEBUG":
		lvl = log.DebugLevel

	case "INFO":
		lvl = log.InfoLevel

	case "WARN":
		lvl = log.WarnLevel

	case "ERROR":
		lvl = log.ErrorLevel

	case "FATAL":
		lvl = log.FatalLevel

	case "PANIC":
		lvl = log.PanicLevel

	default:
		err = fmt.Errorf("invalid error level '%s'", lvlStr)
	}

	return
}
