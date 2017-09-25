package sdk

import (
	"fmt"
	"errors"
	"encoding/json"
	"github.com/op/go-logging"
	"github.com/streadway/amqp"
	"github.com/openbaton/go-openbaton/catalogue"
)

type Handler interface{}

type handleVnfmFunction func(bytemsg []byte, worker interface{}) ([]byte, error)

type handlePluginFunction func(bytemsg []byte, worker interface{}) ([]byte, error)

type manager interface {
	Start(confPath string, name *string)
	getCreds() catalogue.ManagerCredentials
}

func GetVnfmCreds(username string, password string, brokerIp string, brokerPort int, vnfm_endpoint *catalogue.Endpoint, log_level string) (*catalogue.ManagerCredentials, error) {
	registerMessage := catalogue.VnfmRegisterMessage{}
	registerMessage.Action = "register"
	registerMessage.Endpoint = vnfm_endpoint
	registerMessage.Type = vnfm_endpoint.Type
	return getCreds(username, password, brokerIp, brokerPort, registerMessage, log_level)
}

func GetPluginCreds(username string, password string, brokerIp string, brokerPort int, plugin_type string, log_level string) (*catalogue.ManagerCredentials, error) {
	registerMessage := catalogue.PluginRegisterMessage{}
	registerMessage.Action = "register"
	registerMessage.Type = plugin_type
	return getCreds(username, password, brokerIp, brokerPort, registerMessage, log_level)
}

func getCreds(username string, password string, brokerIp string, brokerPort int, msg interface{}, log_level string) (*catalogue.ManagerCredentials, error) {
	amqpUri := getAmqpUri(username, password, brokerIp, brokerPort)
	logger := GetLogger("common.sdk", log_level)
	logger.Debugf("Dialing %s", amqpUri)

	conn, err := amqp.Dial(amqpUri)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	defer channel.Close()
	q, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}

	corrId := RandomString(32)
	mrs, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("Error while marshaling: %v", err)
		return nil, err
	}
	err = channel.Publish(
		"",                      // exchange
		"nfvo.manager.handling", // routing key
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(mrs),
		})
	if err != nil {
		return nil, err
	}

	for d := range msgs {
		if corrId == d.CorrelationId {
			managerCredentials := &catalogue.ManagerCredentials{}
			err := json.Unmarshal(d.Body, managerCredentials)
			if err != nil {
				return nil, err
			}
			return managerCredentials, nil
		}
	}
	return nil, errors.New("no answer")
}

func getAmqpUri(username string, password string, brokerIp string, brokerPort int) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, brokerIp, brokerPort)
}

type Manager struct {
	conn       *amqp.Connection
	Channel    *amqp.Channel
	workers    int
	queueName  string
	done       chan error
	logger     *logging.Logger
	deliveries <-chan amqp.Delivery
}

type VnfmManager struct {
	*Manager
	handlerFunction handleVnfmFunction
}

type PluginManager struct {
	*Manager
	handlerFunction handlePluginFunction
}

func NewPluginManager(username string,
	password string,
	brokerIp string,
	brokerPort int,
	exchange string,
	queueName string,
	workers int,
	handlerFunction handlePluginFunction,
	log_level string) (*PluginManager, error) {
	m := &PluginManager{
		Manager: &Manager{
			logger:  GetLogger("plugin-manager", log_level),
			Channel: nil,
			conn:    nil,
			workers: workers,
			done:    make(chan error),
		},
		handlerFunction: handlerFunction,
	}
	err := setupManager(username, password, brokerIp, brokerPort, m.Manager, exchange, queueName)
	m.queueName = queueName
	if err != nil {
		m.logger.Errorf("Error while setup the amqp thing: %v", err)
		return nil, err
	}

	return m, nil
}

func NewVnfmManager(username string,
	password string,
	brokerIp string,
	brokerPort int,
	exchange string,
	queueName string,
	workers int,
	manager_name string,
	handleFunction handleVnfmFunction,
	log_level string) (*VnfmManager, error) {

	c := &VnfmManager{
		Manager: &Manager{
			conn:    nil,
			Channel: nil,
			workers: workers,
			done:    make(chan error),
			logger:  GetLogger(manager_name, log_level),
		},
		handlerFunction: handleFunction,
	}

	var err error
	err = setupManager(username, password, brokerIp, brokerPort, c.Manager, exchange, queueName)
	if err != nil {
		c.logger.Errorf("Error while setup the amqp thing: %v", err)
		return nil, err
	}
	c.queueName = queueName
	return c, nil
}

func setupManager(username string, password string, brokerIp string, brokerPort int, c *Manager, exchange string, queueName string) error {
	amqpURI := getAmqpUri(username, password, brokerIp, brokerPort)
	c.logger.Debugf("dialing %s", amqpURI)
	var err error
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return err
	}

	go func() {
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	c.logger.Debugf("got Connection, getting Channel")
	c.Channel, err = c.conn.Channel()
	if err != nil {
		return err
	}

	c.logger.Debugf("got Channel, declaring Exchange (%q)", exchange)

	c.logger.Debugf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.Channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		true,      // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	c.logger.Debugf("declared Queue (%q, %d messages, %d consumers), binding to Exchange",
		queue.Name, queue.Messages, queue.Consumers)

	if err = c.Channel.QueueBind(
		queue.Name, // name of the queue
		queue.Name, // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return err
	}

	c.logger.Debug("Queue bound to Exchange, starting Consume")
	return nil
}

func (c *Manager) Shutdown() error {
	if err := c.conn.Close(); err != nil {
		c.logger.Errorf("AMQP connection close error: %s", err)
		return err
	}

	defer c.logger.Debugf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func (v *Manager) Unregister(typ string) {
	msg := catalogue.PluginRegisterMessage{
		Type:typ,
		Action:"unregister",
	}
	resp, _, err := ExecuteRpc("nfvo.manager.handling",msg, v.Channel, v.logger)
	if err != nil {
		v.logger.Errorf("Error unregistering: %v", err)
		return
	}
	v.logger.Debugf("Unregistered and got answer: %v", resp)
}

func (c *VnfmManager) Serve(worker interface{}) {
	forever := make(chan bool)

	for x := 0; x < c.workers; x++ {

		go func() {

			deliveries, err := c.Channel.Consume(
				c.queueName, // name
				"",          // consumerTag,
				false,       // noAck
				false,       // exclusive
				false,       // noLocal
				false,       // noWait
				nil,         // arguments
			)
			if err != nil {
				c.logger.Errorf("Error while consuming: %v", err)
				return
			}

			c.deliveries = deliveries
			for d := range c.deliveries {

				byteRes, err := c.handlerFunction(d.Body, worker)
				if err != nil {
					c.logger.Errorf("Error while executing handler function: %v", err)
					return
				}
				rerr := c.Channel.Publish(
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          byteRes,
					})
				if err != nil {
					c.done <- rerr
					return
				}

				d.Ack(false)
			}
		}()
	}
	<-forever
}

func (c *PluginManager) Serve(worker interface{}) {
	forever := make(chan bool)

	for x := 0; x < c.workers; x++ {

		go func() {

			deliveries, err := c.Channel.Consume(
				c.queueName, // name
				"",          // consumerTag,
				false,       // noAck
				false,       // exclusive
				false,       // noLocal
				false,       // noWait
				nil,         // arguments
			)
			if err != nil {
				c.logger.Errorf("Error while consuming: %v", err)
				return
			}

			c.deliveries = deliveries
			for d := range c.deliveries {

				byteRes, err := c.handlerFunction(d.Body, worker)
				if err != nil {
					c.logger.Errorf("Error while executing handler function: %v", err)
					return
				}
				rerr := c.Channel.Publish(
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          byteRes,
					})
				if err != nil {
					c.done <- rerr
					return
				}

				d.Ack(false)
			}
		}()
	}
	<-forever
}
