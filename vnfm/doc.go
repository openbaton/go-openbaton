/*
Package vnfm implements an extensible, transport-agnostic runtime for OpenBaton VNFMs.

See the go-dummy-vnfm repo for a sample implementation of a VNFM using this package and the AMQP driver.

vnfm uses the vnfm/channel package to abstract the underlying transport channel.
The required drivers must be registered before creating a new VNFM using vnfm.Register(); usually, this is done automatically by the driver package when first imported.

A new VNFM can be created by using the New() function together with a Handler instance:


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

The new VNFM can then be started using its Serve() method, blocking the current goroutine.
Use Stop() to stop the service and quit.

    if err := svc.Serve(); err != nil {
        panic("error while setting up plugin: " + err.Error())
    }
*/
package vnfm
