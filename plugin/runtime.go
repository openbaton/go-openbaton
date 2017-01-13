package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/mcilloni/go-openbaton/util"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const pluginExchange = "plugin-exchange"

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

	// Type is a string that identifies the type of this plugin.
	Type string

	// Workers determines how many workers the plugin will spawn.
	// Set this number according to your needs.
	Workers int

	// LogLevel sets the minimum logging level for the internal instance of logrus.Logger.
	LogLevel log.Level
}

//  Plugin represents a plugin instance.
type Plugin interface {
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

// New creates a plugin from an implementation and plugin.Params.
// impl must be of a valid Plugin implementation type, like plugin.Driver.
func New(impl interface{}, p *Params) (Plugin, error) {
	tag := util.FuncName()

	if p.Workers < 1 {
		p.Workers = 10
	}

	plug := &plug{
		connstr:              util.AmqpUriBuilder(p.Username, p.Password, p.BrokerAddress, "", p.Port, false),
		params:               p,
		quitChan:             make(chan error),
		receiverDeliveryChan: make(chan (<-chan amqp.Delivery), 1),
		reqChan:              make(chan request, 30),
	}

	if err := plug.initLogger(); err != nil {
		return nil, err
	}

	var rh reqHandler

	switch v := impl.(type) {
	case Driver:
		rh = driverHandler{v, plug.l}

	// in case we are reinitialising the plugin
	case reqHandler:
		rh = v

	default:
		plug.l.WithField("tag", tag).Panicf("invalid plugin implementation %T", impl)
	}

	plug.rh = rh

	return plug, nil
}

type plug struct {
	cnl     *amqp.Channel
	conn    *amqp.Connection
	connstr string

	l                    *log.Logger
	e                    logData
	params               *Params
	quitChan             chan error
	receiverDeliveryChan chan (<-chan amqp.Delivery)
	reqChan              chan request
	rh                   reqHandler
	stopped              bool
	wg                   sync.WaitGroup
}

func (p *plug) Logger() *log.Logger {
	return p.l
}

func (p *plug) Serve() error {
	tag := util.FuncName()

	// reinit the plugin if already stopped
	if p.stopped {
		panic("plugin already stopped")
	}

	p.l.WithFields(log.Fields{
		"tag":    tag,
		"params": *p.params,
	}).Debug("plugin starting")

	errChan, err := p.setup()
	if err != nil {
		return err
	}

	p.spawnWorkers()
	p.spawnReceiver()

MainLoop:
	for {
		select {
		case <-p.quitChan:
			if err = p.conn.Close(); err != nil {
				p.l.WithError(err).WithFields(log.Fields{
					"tag": tag,
				}).Error("closing Connection failed")

				p.closeQueues()

				// send the error to stop
				p.quitChan <- err
				return nil
			}

			p.l.WithFields(log.Fields{
				"tag": tag,
			}).Info("initiating clean shutdown")

			// Close will cause the reception of nil on errChan.

		case amqpErr := <-errChan:
			// The connection closed cleanly after invoking Close().
			if amqpErr == nil {
				// notify the receiver and workers
				p.closeQueues()

				p.wg.Wait()

				// send nil to Stop
				close(p.quitChan)

				p.l.WithFields(log.Fields{
					"tag": tag,
				}).Debug("main loop quitting")

				break MainLoop
			}

			p.l.WithError(amqpErr).WithFields(log.Fields{
				"tag": tag,
			}).Error("received AMQP error for current connection")

			// The connection crashed for some reason. Try to bring it up again.
			for {
				if errChan, err = p.setup(); err != nil {
					p.l.WithError(err).WithFields(log.Fields{
						"tag": tag,
					}).Error("can't re-establish connection with AMQP; queues stalled. Retrying in 30 seconds.")
					time.Sleep(30 * time.Second)
				}
			}

		}
	}

	return nil
}

func (p *plug) Stop() error {
	tag := util.FuncName()

	defer p.deinitLogger()

	if p.stopped {
		return fmt.Errorf("plugin %s already stopped", p.params.Name)
	}

	// first step: signal the main routine to quit.
	select {
	case p.quitChan <- nil:

	case <-time.After(time.Second):
		return errors.New("the plugin is not listening")
	}

	// second step: wait for it to quit
	select {
	case err := <-p.quitChan:
		if err != nil {
			return err
		}
	case <-time.After(1 * time.Minute):
		return errors.New("the plugin refused to quit")
	}

	p.stopped = true

	p.l.WithFields(log.Fields{
		"tag": tag,
	}).Info("plugin stopped cleanly")

	return nil
}

func (p *plug) Type() string {
	return p.rh.Type()
}

func (p *plug) closeQueues() {
	// closes the workers
	close(p.reqChan)

	// closes the receiver
	close(p.receiverDeliveryChan)
}

func (p *plug) id() string {
	return fmt.Sprintf("%s.%s.%s", p.rh.QueueTag(), p.params.Type, p.params.Name)
}

func (p *plug) setup() (<-chan *amqp.Error, error) {
	tag := util.FuncName()

	p.l.WithFields(log.Fields{
		"tag": tag,
	}).Info("dialing AMQP")

	conn, err := amqp.Dial(p.connstr)
	if err != nil {
		return nil, err
	}

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

	p.conn = conn
	p.cnl = cnl

	// setup incoming deliveries
	deliveries, err := cnl.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, err
	}

	p.receiverDeliveryChan <- deliveries

	return conn.NotifyClose(make(chan *amqp.Error)), nil
}

type reqHandler interface {
	Handle(call string, args []json.RawMessage) (interface{}, error)
	QueueTag() string
	Type() string
}
