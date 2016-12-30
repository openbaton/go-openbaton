package amqp

import (
	"bytes"
	"strconv"

	"github.com/mcilloni/go-openbaton/vnfm"
	"github.com/mcilloni/go-openbaton/vnfm/channel"
	"github.com/mcilloni/go-openbaton/vnfm/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	vnfm.Register("amqp", amqpDriver{})
}

type amqpDriver struct{}

func (amqpDriver) Init(cnf *config.Config, log *log.Logger) (channel.Channel, error) {
	ret, err := newChannel(cnf, log)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func uriBuilder(username, password, host, vhost string, port int, tls bool) string {
	buffer := bytes.NewBufferString("amqp")

	if tls {
		buffer.WriteRune('s')
	}

	buffer.WriteString("://")
	if username != "" {
		buffer.WriteString(username)

		if password != "" {
			buffer.WriteRune(':')
			buffer.WriteString(password)
		}

		buffer.WriteRune('@')
	}

	if host != "" {
		buffer.WriteString(host)
	}

	if port > 0 {
		buffer.WriteRune(':')
		buffer.WriteString(strconv.Itoa(port))
	}

	if vhost != "" {
		buffer.WriteRune('/')
		buffer.WriteString(vhost)
	}

	return buffer.String()
}
