/*
Package channel abstract the transport channel between an OpenBaton VNFM and the NFVO.

See go-dummy-vnfm for a sample implementation of a VNFM using the AMQP driver.

vnfm uses the vnfm/channel package to abstract the underlying transport channel.
The required drivers must be registered before creating a new VNFM using vnfm.Register(); usually, this is done automatically by the driver package when first imported.

	// import the driver
	import _ "driver/package/xyz"

	//  some code here

	// "xyz" is the identifier of the desired driver.
	svc, err := vnfm.New("xyz", handler, cfg)
	// use the svc
*/
package channel

import (
	"github.com/mcilloni/go-openbaton/catalogue/messages"
	"github.com/mcilloni/go-openbaton/vnfm/config"
	log "github.com/sirupsen/logrus"
)

// Status represents the current status of the channel.
type Status int

const (
	// Running indicates that the channel is connected and working correctly.
	Running Status = iota

	// Reconnecting indicates that the channel is trying to reestablish the connection with the NFVO.
	Reconnecting

	// Stopped indicates that the channel is not active.
	Stopped

	// Quitting indicates that the channel is quitting.
	Quitting
)

// Driver is an interface representing a driver implementation for a channel type.
// A new channel of the associated type can be obtained by using Init().
type Driver interface {
	// Init initialises a Channel instance using the given config.Config.
	// conf.Properties must contain all the values required by the current implementation.
	Init(conf *config.Config, log *log.Logger) (Channel, error)
}

// Channel is an interface that abstracts the transport layer between the NFVO and the VNFM.
type Channel interface {

	// Close closes the channel.
	Close() error

	// Exchange sends a message to an implementation defined destination, and then waits for a reply.
	Exchange(dest string, msg []byte) ([]byte, error)

	// Impl returns the underlying implementation of this channel (if any), i.e., an amqp Channel will return the
	// AMQP Channel underlying.
	// The caller must know if the channel is of the correct type.
	Impl() (interface{}, error)

	// NFVOExchange sends a message to the NFVO, and then waits for a reply.
	// The outgoing message must have From() == messages.VNFR.
	NFVOExchange(msg messages.NFVMessage) (messages.NFVMessage, error)

	// NFVOSend sends a message to the NFVO without waiting for a reply.
	// A success while sending the message is no guarantee about the NFVO actually receiving it.
	NFVOSend(msg messages.NFVMessage) error

	// NotifyReceiver creates a channel on which received messages will be delivered.
	// The returned channel will be removed if nobody is listening on it for a while.
	NotifyReceived() (<-chan messages.NFVMessage, error)

	// Send sends a message to an implementation defined destination without waiting for a reply.
	// A success while sending the message is no guarantee about the destination actually receiving it.
	Send(dest string, msg []byte) error

	// Status returns the current status of the Channel.
	Status() Status
}

// NFVOResponse represents a response from the NFVO.
type NFVOResponse struct {
	messages.NFVMessage
	error
}

// NFVOExchangeAsync executes the NFVOExchange method of a given channel asyncronously, returning
// a channel on which the response will be delivered.
func NFVOExchangeAsync(cnl Channel, msg messages.NFVMessage) <-chan *NFVOResponse {
	ret := make(chan *NFVOResponse, 1)

	go func() {
		msg, err := cnl.NFVOExchange(msg)

		ret <- &NFVOResponse{msg, err}
	}()

	return ret
}
