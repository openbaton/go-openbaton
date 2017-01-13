# Plugin

`plugin` implements a runtime for OpenBaton plugins. 

Currently, only vim-drivers are supported; see [go-vimdriver-test] for a sample implementation.

## Implementing a VIM driver

A new plugin can be created by using the `plugin.New(interface{}, *plugin.Params) (plugin.Plugin, error)` function together with a `plugin.Driver` instance:

```go
var driver plugin.Driver = &myDriver{}

params := &plugin.Params{ /* your configuration here */ }

plug, err := plugin.New(driver, params)
if err != nil {
    panic("error: " + err.Error())
}
```

(Ensure that your VIMDriver implements the `plugin.Driver` interface!)

The new `plugin.Plugin` can then be started using its `Serve()` method, blocking the current goroutine.
Use `Stop()` to stop the service and quit.

```go
if err := plug.Serve(); err != nil {
    panic("error while setting up plugin: " + err.Error())
}
```

## For further informations

Check out the [GoDoc][godoc].

[godoc]: http://godoc.org/github.com/mcilloni/go-openbaton/plugin
[go-vimdriver-test]: https://github.com/mcilloni/go-vimdriver-test