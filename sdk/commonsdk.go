/*
	Common package for the VNFM and the Plugin SDK for Open Baton Managers
 */
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

// Handler function for the VNFMs
type handlerFunction func(bytemsg []byte, handlerVnfm Handler, allocate bool, connection *amqp.Connection, net catalogue.BaseNetworkInt, img catalogue.BaseImageInt) ([]byte, error)

// Function to retrieve the private amqp credentials for a VNFM
func GetVnfmCreds(username string, password string, brokerIp string, brokerPort int, vnfm_endpoint *catalogue.Endpoint, log_level string) (*catalogue.ManagerCredentials, error) {
	registerMessage := catalogue.VnfmRegisterMessage{}
	registerMessage.Action = "register"
	registerMessage.Endpoint = vnfm_endpoint
	registerMessage.Type = vnfm_endpoint.Type
	return getCreds(username, password, brokerIp, brokerPort, registerMessage, log_level)
}

// Function to retrieve the private amqp credentials for a Plugin
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

	conn, err := amqp.DialConfig(amqpUri, amqp.Config{
		Heartbeat: 5,
	})
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

	corrId := randomString(32)
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

// Base Manager struct
type Manager struct {
	Connection      *amqp.Connection
	Channel         *amqp.Channel
	workers         int
	allocate        bool
	queueName       string
	errorChan       chan error
	logger          *logging.Logger
	deliveries      <-chan amqp.Delivery
	handlerFunction handlerFunction
	handler         Handler
	image           catalogue.BaseImageInt
	network         catalogue.BaseNetworkInt
}

// Instantiate a new Manager struct
func NewManager(h Handler,
	username string,
	password string,
	brokerIp string,
	brokerPort int,
	exchange string,
	queueName string,
	workers int,
	allocate bool,
	managerName string,
	handleFunction handlerFunction,
	logLevel string,
	net catalogue.BaseNetworkInt,
	img catalogue.BaseImageInt) (*Manager, error) {

	manager := &Manager{
		Connection:      nil,
		Channel:         nil,
		allocate:        allocate,
		workers:         workers,
		errorChan:       make(chan error),
		logger:          GetLogger(managerName, logLevel),
		handlerFunction: handleFunction,
		handler:         h,
		image:           img,
		network:         net,
	}

	err := setupManager(username, password, brokerIp, brokerPort, manager, exchange, queueName)
	if err != nil {
		manager.logger.Errorf("Error while setup the amqp thing: %v", err)
		return nil, err
	}
	manager.queueName = queueName
	return manager, nil
}

func setupManager(username string, password string, brokerIp string, brokerPort int, c *Manager, exchange string, queueName string) error {
	amqpURI := getAmqpUri(username, password, brokerIp, brokerPort)
	c.logger.Debugf("dialing %s", amqpURI)
	var err error
	c.Connection, err = amqp.Dial(amqpURI)
	if err != nil {
		return err
	}

	c.logger.Debugf("got Connection, getting Channel")
	c.Channel, err = c.Connection.Channel()
	if err != nil {
		return err
	}

	c.logger.Debugf("got Channel, declaring Exchange (%q)", exchange)

	c.logger.Debugf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.Channel.QueueDeclare(
		queueName,
		true,
		true,
		false,
		false,
		nil,
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

// Shutdown the manager
func (manager *Manager) Shutdown() error {
	if err := manager.Connection.Close(); err != nil {
		manager.logger.Errorf("AMQP connection close error: %s", err)
		return err
	}

	defer manager.logger.Debugf("AMQP shutdown OK")

	return <-manager.errorChan
}

// Unregister function for Managers
func (manager *Manager) Unregister(typ, username, password string, vnfmEndpoint *catalogue.Endpoint) {
	if vnfmEndpoint == nil {
		manager.unregisterPlugin(typ, username, password)
		return
	} else {
		msg := catalogue.VnfmManagerUnregisterMessage{
			Type:     typ,
			Action:   "unregister",
			Username: username,
			Password: password,
			Endpoint: vnfmEndpoint,
		}
		manager.unregister(msg)
	}
}

// Unregister function for the Plugin
func (manager *Manager) unregisterPlugin(typ, username, password string) {
	msg := catalogue.ManagerUnregisterMessage{
		Type:     typ,
		Action:   "unregister",
		Username: username,
		Password: password,
	}
	manager.unregister(msg)
}

func (manager *Manager) unregister(msg interface{}) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		manager.logger.Errorf("Error while marshalling unregister message: %v", err)
		return
	}
	err = SendMsg("nfvo.manager.handling", msgBytes, manager.Channel, manager.logger)
	if err != nil {
		manager.logger.Errorf("Error unregistering: %v", err)
		return
	}
}

// Serve function for Manager
func (manager *Manager) Serve() {
	forever := make(chan bool)

	for x := 0; x < manager.workers; x++ {

		go func() {

			deliveries, err := manager.Channel.Consume(
				manager.queueName,
				"",
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				manager.logger.Errorf("Error while consuming: %v", err)
				return
			}

			manager.deliveries = deliveries
			for d := range manager.deliveries {
				d1 := d
				go func() {
					byteRes, err := manager.handlerFunction(d1.Body, manager.handler, manager.allocate, manager.Connection, manager.network, manager.image)
					if err != nil {
						manager.logger.Errorf("Error while executing handler function: %v", err)
						return
					}
					err = manager.Channel.Publish(
						"",
						d1.ReplyTo,
						false,
						false,
						amqp.Publishing{
							ContentType:   "text/plain",
							CorrelationId: d1.CorrelationId,
							Body:          byteRes,
						})
					if err != nil {
						manager.errorChan <- err
						return
					}
				}()

				d.Ack(false)
			}
		}()
	}
	go func() {
		for {
			manager.logger.Error(fmt.Sprintf("Got error while handling rabbitmq: %q", <-manager.errorChan))
		}
	}()
	<-forever
}
