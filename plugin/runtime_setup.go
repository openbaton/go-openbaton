package plugin

import (
	"time"

	"github.com/openbaton/go-openbaton/util"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func (p *plug) connSetup() (*amqp.Connection, error) {
	tag := util.FuncName()

	p.l.WithFields(log.Fields{
		"tag": tag,
	}).Info("dialing AMQP")

	return amqp.Dial(p.connstr)
}

func (p *plug) getAMQPChan() (*amqp.Channel, error) {
	select {
	case p.chanReqChan <- struct{}{}:
		// ok, sent

	case <-time.After(5 * time.Second):
		return nil, ErrTimeout
	}

	nc, ok := <-p.newChanChan
	if !ok {
		return nil, nil
	}

	return nc.Channel, nc.error
}

func (p *plug) makeAMQPChan(conn *amqp.Connection) (*amqp.Channel, error) {
	cnl, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := cnl.ExchangeDeclare(pluginExchange, "topic", false,
		false, false, false, nil); err != nil {
		return nil, err
	}

	queueName := p.id()
	if _, err := cnl.QueueDeclare(queueName, false, true,
		false, false, nil); err != nil {
		return nil, err
	}

	if err := cnl.QueueBind(queueName, queueName, pluginExchange, false, nil); err != nil {
		return nil, err
	}

	if err := cnl.Qos(1, 0, false); err != nil {
		return nil, err
	}

	return cnl, nil
}
