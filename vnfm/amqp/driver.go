package amqp

import (
	"github.com/mcilloni/go-openbaton/vnfm"
	"github.com/mcilloni/go-openbaton/vnfm/channel"
	"github.com/mcilloni/go-openbaton/vnfm/config"
)

func init() {
	vnfm.Register("amqp", amqpDriver{})
}

type amqpDriver struct{}

func (amqpDriver) Init(props config.Properties) (channel.Channel, error) {
	return newChannel()
}
