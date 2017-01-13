# vnfm

`vmfm` implements an extensible, transport-agnostic runtime for OpenBaton VNFMs. 

See [go-dummy-vnfm] for a sample implementation of a VNFM using this package and the AMQP driver.

## Transports

`vnfm` uses the `vnfm/channel` package to abstract the underlying transport channel.
The required drivers must be registered before creating a new VNFM using `vnfm.Register()`; usually, this is done automatically by the driver package when first imported.

## Implementing a VNFM

A new VNFM can be created by using the `vnfm.New(string, vnfm.Handler, *config.Config) (vnfm.VNFM, error)` function together with a `vnfm.Handler` instance:

```go
// import the driver
import _ "driver/package/xyz"

var handler vnfm.Handler = &myHandler{}

cfg, err := config.LoadFile("path/to/config.toml")
if err != nil {
    panic("cannot load config, " + err.Error())
}

// "xyz" is the identifier of the desired driver.
svc, err := vnfm.New("xyz", handler, cfg)
if err != nil {
    panic("error: " + err.Error())
}
```

(Ensure that your handler implements the `vnfm.Handler` interface!)

The new `vnfm.Handler` can then be started using its `Serve()` method, blocking the current goroutine.
Use `Stop()` to stop the service and quit.

```go
if err := svc.Serve(); err != nil {
    panic("error while setting up plugin: " + err.Error())
}
```

## For further informations

Check out the [GoDoc][godoc].

[godoc]: http://godoc.org/github.com/mcilloni/go-openbaton/vnfm
[go-dummy-vnfm]: https://github.com/mcilloni/go-dummy-vnfm