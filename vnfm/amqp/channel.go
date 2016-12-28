package amqp

import (
	"errors"
	"time"

	"github.com/mcilloni/go-openbaton/catalogue/messages"
	"github.com/mcilloni/go-openbaton/log"
	"github.com/mcilloni/go-openbaton/vnfm/config"
	"github.com/streadway/amqp"
	"github.com/mcilloni/go-openbaton/vnfm/channel"
)

type amqpConf struct {
	connstr string
	cfg     amqp.Config
}

type exchange struct {
	payload
	replyChan chan messages.NFVMessage
}

type payload struct {
}

// An amqpChannel is a control structure to handle an AMQP connection.
// The main logic is handled in an event loop, which is fed using Go channels through
// the amqpChannel methods.
type amqpChannel struct {
	acfg         amqpConf
	exchangeChan chan *exchange
	l            *log.Logger
	notifyQueue  []chan<- messages.NFVMessage
	quitChan     chan struct{}
	sendChan     chan *payload
	status		 channel.Status
	subChan   	 chan chan messages.NFVMessage
}

func newChannel(log *log.Logger) *amqpChannel {
	return &amqpChannel{
		exchangeChan: make(chan *exchange),
		l:            log,
		notifyQueue:  []chan<- messages.NFVMessage{},
		quitChan:     make(chan struct{}),
		sendChan:     make(chan *payload, 20),
		status:		  channel.Stopped,
		subChan:	  make(chan chan messages.NFVMessage),
	}
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

func (acnl *amqpChannel) Exchange(msg messages.NFVMessage, timeout time.Duration) (messages.NFVMessage, error) {
	return nil, nil
}

func (acnl *amqpChannel) ExchangeStrings(msg, queue string, timeout time.Duration) (string, error) {
	return "", nil
}

func (acnl *amqpChannel) NotifyReceived() (<-chan messages.NFVMessage, error) {
	notifyChan := make(chan messages.NFVMessage, 5)

	acnl.subChan <- notifyChan

	return notifyChan, nil
}

func (acnl *amqpChannel) Send(msg messages.NFVMessage) error {
	return nil
}

func (acnl *amqpChannel) Status() channel.Status {
	return acnl.status
}

func (acnl *amqpChannel) broadcastNotification(msg messages.NFVMessage) {
	newList := make([]chan<- messages.NFVMessage, len(acnl.notifyQueue))

	for _, c := range acnl.notifyQueue {
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

	acnl.notifyQueue = newList
}

func (acnl *amqpChannel) closeQueues() {
	close(acnl.quitChan)

	for _, cnl := range acnl.notifyQueue {
		close(cnl)
	}
}

func (acnl *amqpChannel) setupWithProps(props config.Properties) error {
	// defaults
	host := "localhost"
	port := 5672
	username := ""
	password := ""
	vhost := ""

	sect, ok := props.Section("amqp")
	if ok {
		acnl.l.Infoln("found AMQP section in config")

		host, _ = sect.ValueString("host", "localhost")
		username, _ = sect.ValueString("username", "")
		password, _ = sect.ValueString("password", "")
		port, _ = sect.ValueInt("port", 5672)
		vhost, _ = sect.ValueString("vhost", "")
	}

	// TODO: handle TLS
	acnl.acfg.connstr = uriBuilder(username, password, host, vhost, port, false)

	heartbeat, _ := sect.ValueInt("heartbeat", 60)

	acnl.acfg.cfg = amqp.Config{
		Heartbeat: time.Duration(heartbeat) * time.Second,
	}

	return nil
}

func (acnl *amqpChannel) setup() (*amqp.Connection, *amqp.Channel, chan *amqp.Error, error) {
	acnl.l.Infof("dialing AMQP with uri %s\n", acnl.acfg.connstr)

	conn, err := amqp.DialConfig(acnl.acfg.connstr, acnl.acfg.cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	cnl, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, err
	}

	acnl.status = channel.Running

	return conn, cnl, conn.NotifyClose(make(chan *amqp.Error)), nil
}

// spawn spawns the main handler for AMQP communications.
func (acnl *amqpChannel) spawn() error {
	conn, cnl, errChan, err := acnl.setup()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case notifyChan := <-acnl.subChan:
				acnl.notifyQueue = append(acnl.notifyQueue, notifyChan)
			
			case <-acnl.quitChan:
				if err := conn.Close(); err != nil {
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

				acnl.status = channel.Reconnecting

				// The connection crashed for some reason. Try to bring it up again.
				for {
					if conn, cnl, errChan, err = acnl.setup(); err != nil {
						acnl.l.Errorln("can't re-establish connection with AMQP; queues stalled. Retrying in 30 seconds.")
						time.Sleep(30 * time.Second)
					} else {
						break
					}
				}

			case payload := <-acnl.sendChan:

			case exchange := <-acnl.exchangeChan:
			}
		}
	}()

	return nil
}
