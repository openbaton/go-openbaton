package amqp

import (
	"errors"
	"time"

	"github.com/mcilloni/go-openbaton/catalogue"
	"github.com/mcilloni/go-openbaton/catalogue/messages"
	"github.com/mcilloni/go-openbaton/log"
	"github.com/mcilloni/go-openbaton/vnfm/channel"
	"github.com/mcilloni/go-openbaton/vnfm/config"
	"github.com/streadway/amqp"
)

var (
	ErrTimedOut = errors.New("timed out")
)

type exchange struct {
	queue     string
	msg       []byte
	replyChan chan response
}

type exchangeTicket struct {
	id       catalogue.ID
	respChan chan<- []byte
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
		connstr  string
		cfg      amqp.Config
		exchange struct {
			name    string
			durable bool
		}
	}

	conn *amqp.Connection
	cnl  *amqp.Channel

	l            *log.Logger
	notifyChans  []chan<- messages.NFVMessage
	numOfWorkers int
	quitChan     chan struct{}
	sendQueue    chan *exchange
	status       channel.Status
	statusChan   chan channel.Status
	subChan      chan chan messages.NFVMessage
}

func newChannel(props config.Properties, log *log.Logger) (*amqpChannel, error) {
	acnl := &amqpChannel{
		l:           log,
		notifyChans: []chan<- messages.NFVMessage{},
		quitChan:    make(chan struct{}),
		status:      channel.Stopped,
		subChan:     make(chan chan messages.NFVMessage),
	}

	// defaults
	host := "localhost"
	port := 5672
	username := ""
	password := ""
	vhost := ""
	heartbeat := 60
	exchangeName := ExchangeDefault
	exchangeDurable := true

	workers, queueSize := 5, 20

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

		workers, _ = sect.ValueInt("workers", workers)
		queueSize, _ = sect.ValueInt("queue_size", queueSize)
	}

	// TODO: handle TLS
	acnl.cfg.connstr = uriBuilder(username, password, host, vhost, port, false)

	acnl.cfg.cfg = amqp.Config{
		Heartbeat: time.Duration(heartbeat) * time.Second,
	}

	acnl.cfg.exchange.name = exchangeName
	acnl.cfg.exchange.durable = exchangeDurable

	acnl.sendQueue = make(chan *exchange, queueSize)
	acnl.numOfWorkers = workers
	acnl.statusChan = make(chan channel.Status, workers)

	return acnl, nil
}

func (acnl *amqpChannel) Close() error {
	acnl.quitChan <- struct{}{}

	select {
	case <-acnl.quitChan:
		return nil

	case <-time.After(2 * time.Second):
		return errors.New("timed out afer waiting for AMQP handler loop to close")
	}
}

func (acnl *amqpChannel) Exchange(queue string, msg []byte) ([]byte, error) {
	respChan := make(chan response)

	acnl.sendQueue <- &exchange{queue, msg, respChan}

	resp := <-respChan
	return resp.msg, resp.error
}

func (acnl *amqpChannel) NFVOExchange(msg messages.NFVMessage) (messages.NFVMessage, error) {
	msgBytes, err := messages.Marshal(msg)
	if err != nil {
		return nil, err
	}

	retBytes, err := acnl.Exchange(QueueVNFMCoreActionsReply, msgBytes)
	if err != nil {
		return nil, err
	}

	return messages.Unmarshal(retBytes)
}

func (acnl *amqpChannel) NFVOSend(msg messages.NFVMessage) error {
	msgBytes, err := messages.Marshal(msg)
	if err != nil {
		return err
	}

	return acnl.Send(QueueVNFMCoreActions, msgBytes)
}

func (acnl *amqpChannel) NotifyReceived() (<-chan messages.NFVMessage, error) {
	notifyChan := make(chan messages.NFVMessage, 5)

	acnl.subChan <- notifyChan

	return notifyChan, nil
}

func (acnl *amqpChannel) Send(queue string, msg []byte) error {
	acnl.sendQueue <- &exchange{queue, msg, nil}

	return nil
}

func (acnl *amqpChannel) Status() channel.Status {
	return acnl.status
}

func (acnl *amqpChannel) broadcastNotification(msg messages.NFVMessage) {
	newList := make([]chan<- messages.NFVMessage, len(acnl.notifyChans))

	for _, c := range acnl.notifyChans {
		select {
		// message sent successfully.
		case c <- msg:
			// keep the channel around for the next time
			newList = append(newList, c)

		// nobody is listening at the other end of the channel.
		case <-time.After(1 * time.Second):
			close(c)
		}
	}

	acnl.notifyChans = newList
}

func (acnl *amqpChannel) closeQueues() {
	close(acnl.sendQueue)
	close(acnl.quitChan)

	for _, cnl := range acnl.notifyChans {
		close(cnl)
	}
}

func (acnl *amqpChannel) receiver() {

}

func (acnl *amqpChannel) setStatus(newStatus channel.Status) {
	for i := 0; i < acnl.numOfWorkers; i++ {
		acnl.statusChan <- newStatus
	}

	acnl.status = newStatus
}

func (acnl *amqpChannel) setup() (chan *amqp.Error, error) {
	acnl.l.Infof("dialing AMQP with uri %s\n", acnl.cfg.connstr)

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

	acnl.conn = conn
	acnl.cnl = cnl

	return conn.NotifyClose(make(chan *amqp.Error)), nil
}

// spawn spawns the main handler for AMQP communications.
func (acnl *amqpChannel) spawn() error {
	errChan, err := acnl.setup()
	if err != nil {
		return err
	}

	acnl.spawnWorkers()
	acnl.setStatus(channel.Running)

	go func() {
		for {
			select {
			case notifyChan := <-acnl.subChan:
				acnl.notifyChans = append(acnl.notifyChans, notifyChan)

			case <-acnl.quitChan:
				if err := acnl.conn.Close(); err != nil {
					acnl.l.Errorf("while closing AMQP Connection: %v\n", err)

					acnl.closeQueues()

					return
				}
				// Close will cause the reception of nil on errChan.

			case err = <-errChan:
				// The connection closed cleanly after invoking Close().
				if err == nil {
					// notify the receiving end and listeners
					acnl.closeQueues()

					return
				}

				acnl.setStatus(channel.Reconnecting)

				// The connection crashed for some reason. Try to bring it up again.
				for {
					if errChan, err = acnl.setup(); err != nil {
						acnl.l.Errorln("can't re-establish connection with AMQP; queues stalled. Retrying in 30 seconds.")
						time.Sleep(30 * time.Second)
					} else {
						acnl.setStatus(channel.Running)
						break
					}
				}

			}
		}
	}()

	return nil
}

func (acnl *amqpChannel) spawnWorkers() {
	for i := 0; i < acnl.numOfWorkers; i++ {
		go acnl.worker(i)
	}
}

func (acnl *amqpChannel) worker(id int) {
	acnl.l.Infof("starting worker %d\n", id)

	status := channel.Stopped

	// explanation: a read on a nil channel will
	// block forever. This lambda ensures that we will accept jobs only
	// when the status is valid.
	work := func() chan *exchange {
		if status == channel.Running {
			return acnl.sendQueue
		}

		return nil
	}

	for {
		select {
		case status = <-acnl.statusChan:
			// Updates the status. If it becomes Running, the next loop will accept incoming jobs again

		case exc := <-work():
			// the sender expects a reply
			if exc.replyChan != nil {

			}
		}
	}

	acnl.l.Infof("quitting worker %d\n", id)
}
