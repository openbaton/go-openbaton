package plugin

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Params is a struct containing the plugin's configuration.
type Params struct {
	// BrokerAddress is the address at which the broker AMQP server can be reached.
	BrokerAddress string

	// Port of the AMQP broker.
	Port int

	// Username, Password for the AMQP broker.
	Username, Password string

	// LogFile contains the path to the log file.
	// Use "" to use defaults, or "-" to use stderr.
	LogFile string

	// Name is a parameter provided by the NFVO, usually "openbaton"
	Name string

	// Timestamps enables timestamps.
	Timestamps bool

	// Type is a string that identifies the type of this plugin.
	Type string

	// Workers determines how many workers the plugin will spawn.
	// Set this number according to your needs.
	Workers int

	// LogLevel sets the minimum logging level for the internal instance of logrus.Logger.
	LogLevel log.Level
}

// Plugin represents a plugin instance.
type Plugin interface {
	// ChannelAccessor returns a closure that returns the underlying *amqp.Channel of this Plugin.
	ChannelAccessor() func() (*amqp.Channel, error)

	// Logger returns the internal logger of this Plugin.
	Logger() *log.Logger

	// Serve spawns the Plugin, blocking the current goroutine.
	// Serve only returns non-nil errors during the initialisation phase.
	// Check the log and the return value of Stop() for runtime and on-closing errors respectively.
	Serve() error

	// Stop() signals the event loop of the plugin to quit, and waits until either it shuts down or
	// it times out.
	Stop() error

	// Type() returns the type of this plugin, as specified by its parameters during construction.
	Type() string
}
