package config

import (
	"time"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"errors"
)

// Config represents a generic config type for a VNFM,
// exporting some basic variables from the '[vnfm]' section of 
// the config file.
type Config struct {
    Allocate bool

    LogFile string

    // Properties contain the raw Properties from which this config has
    // been extracted. They also may contain implementation specific fields that 
    // may be needed.
    Properties Properties

    // Timeout represents the amount of time to be waited before timing out.
    Timeout time.Duration
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

	timeoutInt, _ := vnfm.ValueInt("timeout", 2000)
    timeout := time.Duration(timeoutInt) * time.Millisecond

    return &Config{
        Allocate: allocate,
        LogFile: logFile,
        Properties: props,
        Timeout: timeout,
    }, nil
}