# config

`config` handles configurations for VNFMs and relative plugins.

See [go-dummy-vnfm] for a sample implementation of a VNFM using this package and the AMQP driver.

## Overview

`vnfm/config` parses a config as defined by a TOML file.
This config file can specify both generical vnfm parameters (under the [vnfm] section) and driver-specific configurations; 
the package will only read the [vnfm] section into a `config.Config` structure, leaving all the other parameters available to be read from a `config.Properties` instance.

The various drivers can then access their own config sections using this properties instance.

## Loading a config

A config can be read from a file using the `config.LoadFile(string) (*config.Config, error)` function:

```go
cfg, err := config.LoadFile("path/to/config.toml")
if err != nil {
    panic("cannot load config, " + err.Error())
}
```

## For further informations

Check out the [GoDoc][godoc].

[godoc]: http://godoc.org/github.com/openbaton/go-openbaton/vnfm/config
[go-dummy-vnfm]: https://github.com/openbaton/go-dummy-vnfm