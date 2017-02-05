# channel

`channel` abstract the transport channel between an Open Baton VNFM and the NFVO.

See [go-dummy-vnfm] for a sample implementation of a VNFM using the AMQP driver.

## Transports

`vnfm` uses the `vnfm/channel` package to abstract the underlying transport channel.
The required drivers must be registered before creating a new VNFM using `vnfm.Register()`; usually, this is done automatically by the driver package when first imported.

```go
// import the driver
import _ "driver/package/xyz"

/* ...below */

// "xyz" is the identifier of the desired driver.
svc, err := vnfm.New("xyz", handler, cfg)
// use the svc
```

## For further informations

Check out the [GoDoc][godoc].

[godoc]: http://godoc.org/github.com/openbaton/go-openbaton/vnfm/channel
[go-dummy-vnfm]: https://github.com/openbaton/go-dummy-vnfm