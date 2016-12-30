package channel

import (
	"github.com/mcilloni/go-openbaton/catalogue/messages"
	"github.com/mcilloni/go-openbaton/vnfm/config"
	"github.com/mcilloni/go-openbaton/log"
)

type Status int

const (
	Running Status = iota
	Reconnecting
	Stopped
	Quitting
)

type Driver interface {
	// Init initialises a Channel instance using the given config.Config.
	// conf.Properties must contain all the values required by the current implementation.
	Init(conf *config.Config, log *log.Logger) (Channel, error)
}

type Channel interface {
	Close() error

	Exchange(dest string, msg []byte) ([]byte, error)
	
	NFVOExchange(msg messages.NFVMessage) (messages.NFVMessage, error)
	NFVOSend(msg messages.NFVMessage) error

	NotifyReceived() (<-chan messages.NFVMessage, error)

	Send(dest string, msg []byte) error

	Status() Status
}

type NFVOResponse struct {
	messages.NFVMessage
	error
}

func NFVOExchangeAsync(cnl Channel, msg messages.NFVMessage) <-chan *NFVOResponse {
	ret := make(chan *NFVOResponse, 1)

	go func() {
		msg, err := cnl.NFVOExchange(msg)

		ret <- &NFVOResponse{msg, err}
	}()

	return ret
}
