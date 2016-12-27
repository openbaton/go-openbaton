package channel

import (
	"time"

	"github.com/mcilloni/go-openbaton/catalogue/messages"
	"github.com/mcilloni/go-openbaton/vnfm/config"
)

type Driver interface {
	// Init initialises a Channel instance using the given config.Properties.
	// props must contain all the values required by the current implementation.
	Init(props config.Properties) (Channel, error)
}

type Channel interface {
	Close() error

	Exchange(msg messages.NFVMessage, timeout time.Duration) (messages.NFVMessage, error)
	ExchangeStrings(msg, queue string, timeout time.Duration) (string, error)

	NotifyReceived() (<-chan messages.NFVMessage, error)

	Send(msg messages.NFVMessage) error
}

type NFVOResponse struct {
	messages.NFVMessage
	error
}

func ExchangeAsync(cnl Channel, msg messages.NFVMessage, timeout time.Duration) <-chan *NFVOResponse {
	ret := make(chan *NFVOResponse, 1)

	go func() {
		msg, err := cnl.Exchange(msg, timeout)

		ret <- &NFVOResponse{msg, err}
	}()

	return ret
}
