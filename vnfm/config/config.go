package config

import (
	"io"
	"os"

	"errors"
	"github.com/BurntSushi/toml"
)

// Config represents a generic config type for a VNFM,
// exporting some basic variables from the '[vnfm]' section of
// the config file.
type Config struct {
	Allocate bool

	LogFile string

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

	if _, err := toml.DecodeReader(reader, props); err != nil {
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

	logFile, set := vnfm.ValueString("logfile-path", "")
	if !set {
		logFile = ""
	}

	vnfmType, set := vnfm.ValueString("type", "")
	if !set {
		return nil, errors.New("no vnfm.type in config")
	}

	endpoint, set := vnfm.ValueString("endpoint", "")
	if !set {
		return nil, errors.New("no vnfm.endpoint in config")
	}

	descr, _ := vnfm.ValueString("description", "")

	return &Config{
		Allocate:    allocate,
		Description: descr,
		Endpoint:    endpoint,
		LogFile:     logFile,
		Properties:  props,
		Type:        vnfmType,
	}, nil
}
