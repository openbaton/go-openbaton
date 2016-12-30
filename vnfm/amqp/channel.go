package amqp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mcilloni/go-openbaton/catalogue"
	"github.com/mcilloni/go-openbaton/catalogue/messages"
	"github.com/mcilloni/go-openbaton/vnfm/channel"
	"github.com/mcilloni/go-openbaton/vnfm/config"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type exchange struct {
	queue     string
	msg       []byte
	replyChan chan response
}

type response struct {
	msg []byte
	error
}

// An amqpChannel is a control structure to handle an AMQP connection.
// The main logic is handled in an event loop, which is fed using Go channels through
// the amqpChannel methods.
type amqpChannel struct {
	cfg struct {
		connstr string
		cfg     amqp.Config

		exchange struct {
			name    string
			durable bool
		}

		queues struct {
			autodelete, exclusive bool

			generic string
		}

		vnfmType, vnfmEndpoint, vnfmDescr string
	}

	conn *amqp.Connection
	cnl  *amqp.Channel

	receiverDeliveryChan chan (<-chan amqp.Delivery)

	l            *log.Logger
	numOfWorkers int
	quitChan     chan struct{}
	sendQueue    chan *exchange
	status       channel.Status
	statusChan   chan channel.Status
	subChan      chan chan messages.NFVMessage
}

func newChannel(config *config.Config, log *log.Logger) (*amqpChannel, error) {
	props := config.Properties

	acnl := &amqpChannel{
		l:                    log,
		quitChan:             make(chan struct{}),
		receiverDeliveryChan: make(chan (<-chan amqp.Delivery), 1),
		status:               channel.Stopped,
		subChan:              make(chan chan messages.NFVMessage),
	}

	acnl.cfg.vnfmDescr = config.Description
	acnl.cfg.vnfmEndpoint = config.Endpoint
	acnl.cfg.vnfmType = config.Type
	acnl.cfg.queues.generic = fmt.Sprintf("nfvo.%s.actions", config.Type)

	// defaults
	host := "localhost"
	port := 5672
	username := ""
	password := ""
	vhost := ""
	heartbeat := 60
	exchangeName := ExchangeDefault
	exchangeDurable := true
	queuesExclusive := false
	queuesAutodelete := true

	workers, jobQueueSize := 5, 20

	if sect, ok := props.Section("amqp"); ok {
		acnl.l.Infoln("found AMQP section in config")

		host, _ = sect.ValueString("host", host)
		username, _ = sect.ValueString("username", username)
		password, _ = sect.ValueString("password", password)
		port, _ = sect.ValueInt("port", port)
		vhost, _ = sect.ValueString("vhost", vhost)
		heartbeat, _ = sect.ValueInt("heartbeat", heartbeat)

		if exc, ok := sect.Section("exchange"); ok {
			exchangeName, _ = exc.ValueString("name", exchangeName)
			exchangeDurable, _ = exc.ValueBool("durable", exchangeDurable)
		}

		if qus, ok := sect.Section("queues"); ok {
			queuesAutodelete, _ = qus.ValueBool("autodelete", queuesAutodelete)
			queuesExclusive, _ = qus.ValueBool("exclusive", queuesExclusive)
		}

		jobQueueSize, _ = sect.ValueInt("jobqueue-size", jobQueueSize)
		workers, _ = sect.ValueInt("workers", workers)
	}

	// TODO: handle TLS
	acnl.cfg.connstr = uriBuilder(username, password, host, vhost, port, false)

	acnl.cfg.cfg = amqp.Config{
		Heartbeat: time.Duration(heartbeat) * time.Second,
	}

	acnl.cfg.exchange.name = exchangeName
	acnl.cfg.exchange.durable = exchangeDurable

	acnl.cfg.queues.autodelete = queuesAutodelete
	acnl.cfg.queues.exclusive = queuesExclusive

	acnl.sendQueue = make(chan *exchange, jobQueueSize)
	acnl.numOfWorkers = workers
	acnl.statusChan = make(chan channel.Status, workers)

	return acnl, acnl.spawn()
}

func (acnl *amqpChannel) endpoint() *catalogue.Endpoint {
	return &catalogue.Endpoint{
		Active:       true,
		Description:  acnl.cfg.vnfmDescr,
		Enabled:      true,
		Endpoint:     acnl.cfg.vnfmEndpoint,
		EndpointType: "RABBIT",
		ID:           catalogue.GenerateID(),
		Type:         acnl.cfg.vnfmType,
	}
}

func (acnl *amqpChannel) register() error {
	msg, err := json.Marshal(acnl.endpoint())
	if err != nil {
		return err
	}

	return acnl.publish(QueueVNFMRegister, msg)
}

func (acnl *amqpChannel) setup() (<-chan *amqp.Error, error) {
	acnl.l.Infof("dialing AMQP with uri %s", acnl.cfg.connstr)

	conn, err := amqp.DialConfig(acnl.cfg.connstr, acnl.cfg.cfg)
	if err != nil {
		return nil, err
	}

	cnl, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := cnl.ExchangeDeclare(acnl.cfg.exchange.name, "topic", acnl.cfg.exchange.durable,
		false, false, false, nil); err != nil {
		return nil, err
	}

	if err := acnl.setupQueues(cnl); err != nil {
		return nil, err
	}

	acnl.conn = conn
	acnl.cnl = cnl

	// setup incoming deliveries
	deliveries, err := cnl.Consume(
		acnl.cfg.queues.generic, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	acnl.receiverDeliveryChan <- deliveries

	return conn.NotifyClose(make(chan *amqp.Error)), nil
}

// Pretty random, should be checked
func (acnl *amqpChannel) setupQueues(cnl *amqp.Channel) error {
	/*if _, err := cnl.QueueDeclare(QueueVNFMRegister, true, acnl.cfg.queues.autodelete,
		acnl.cfg.queues.exclusive, false, nil); err != nil {

		return err
	}

	if err := cnl.QueueBind(QueueVNFMRegister, QueueVNFMRegister, acnl.cfg.exchange.name, false, nil); err != nil {
		return err
	}

	if _, err := cnl.QueueDeclare(QueueVNFMUnregister, true, acnl.cfg.queues.autodelete,
		acnl.cfg.queues.exclusive, false, nil); err != nil {

		return err
	}

	if err := cnl.QueueBind(QueueVNFMUnregister, QueueVNFMUnregister, acnl.cfg.exchange.name, false, nil); err != nil {
		return err
	}

	if _, err := cnl.QueueDeclare(QueueVNFMCoreActions, true, acnl.cfg.queues.autodelete,
		acnl.cfg.queues.exclusive, false, nil); err != nil {

		return err
	}

	if err := cnl.QueueBind(QueueVNFMCoreActions, QueueVNFMCoreActions, acnl.cfg.exchange.name, false, nil); err != nil {
		return err
	}

	if _, err := cnl.QueueDeclare(QueueVNFMCoreActionsReply, true, acnl.cfg.queues.autodelete,
		acnl.cfg.queues.exclusive, false, nil); err != nil {

		return err
	}

	if err := cnl.QueueBind(QueueVNFMCoreActionsReply, QueueVNFMCoreActionsReply, acnl.cfg.exchange.name, false, nil); err != nil {
		return err
	}*/

	// is this needed?
	if _, err := cnl.QueueDeclare(acnl.cfg.queues.generic, true, acnl.cfg.queues.autodelete,
		acnl.cfg.queues.exclusive, false, nil); err != nil {

		return err
	}

	if err := cnl.QueueBind(acnl.cfg.queues.generic, acnl.cfg.queues.generic, acnl.cfg.exchange.name, false, nil); err != nil {
		return err
	}
	return nil
}

// unregister attempts several times to unregister the Endpoint,
// reestablishing the connection in case of previous failure.
func (acnl *amqpChannel) unregister() error {
	const Attempts = 2

	msg, err := json.Marshal(acnl.endpoint())
	if err != nil {
		return err
	}

	unregFn := func() error {
		return acnl.publish(QueueVNFMUnregister, msg)
	}

	for i := 0; i < Attempts; i++ {
		if i > 0 {
			acnl.l.Warnf("endpoint unregister request failed to send. Reinitializing the connection (tentative #%d)", i)
			if _, err = acnl.setup(); err != nil {
				continue
			}
		}

		if err = unregFn(); err == nil {
			acnl.l.Infof("endpoint unregister request successfully sent at tentative %d", i)
			return nil
		}
	}

	return err
}
