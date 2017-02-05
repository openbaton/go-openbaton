package amqp

import (
	"time"

	"github.com/openbaton/go-openbaton/catalogue/messages"
	"github.com/openbaton/go-openbaton/util"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func (acnl *Channel) handleDelivery(delivery amqp.Delivery, notifyChans []chan<- messages.NFVMessage) ([]chan<- messages.NFVMessage, error) {
	tag := util.FuncName()

	msg, err := messages.Unmarshal(delivery.Body, messages.NFVO)
	if err != nil {
		return nil, err
	}

	acnl.l.WithFields(log.Fields{
		"tag": tag,
		"msg": msg,
	}).Debug("received message")

	last := 0
	for _, c := range notifyChans {
		select {
		// message sent successfully.
		case c <- msg:
			// keep the channel around for the next time
			notifyChans[last] = c
			last++

		// nobody is listening at the other end of the channel.
		case <-time.After(1 * time.Second):
			acnl.l.WithFields(log.Fields{
				"tag": tag,
			}).Debug("closing unresponsive notify channel")

			close(c)
		}
	}

	acnl.l.WithFields(log.Fields{
		"tag":          tag,
		"msg":          msg,
		"num-of-chans": last,
	}).Debug("message dispatched")

	// notifyChans trimmed of dead chans
	return notifyChans[:last], nil
}

func (acnl *Channel) setupDeliveries(cnl *amqp.Channel) (<-chan amqp.Delivery, error) {
	// setup incoming deliveries
	return cnl.Consume(
		acnl.cfg.queues.generic, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

func (acnl *Channel) receiver() {
	tag := util.FuncName()

	acnl.l.WithFields(log.Fields{
		"tag": tag,
	}).Info("AMQP receiver starting")

	// list of channels to which incoming messages will be broadcasted.
	notifyChans := []chan<- messages.NFVMessage{}

RecvLoop:
	for {
		cnl, err := acnl.getAMQPChan()
		if err != nil {
			acnl.l.WithError(err).WithFields(log.Fields{
				"tag": tag,
			}).Error("error while getting an AMQP channel")

			// retry
			continue RecvLoop
		}

		if cnl == nil {
			break RecvLoop
		}

		errChan := cnl.NotifyClose(make(chan *amqp.Error))

		acnl.l.WithFields(log.Fields{
			"tag": tag,
		}).Debug("new AMQP channel received")

		deliveries, err := acnl.setupDeliveries(cnl)
		if err != nil {
			acnl.l.WithError(err).WithFields(log.Fields{
				"tag": tag,
			}).Error("error during delivery handling")
		}

		for {
			select {
			case <-errChan:
				// receving something on errChan means that the channel has closed,
				// either by an issue or because the connection has been closed and
				// the service is closing.
				// In either case, getAMQPChan will return either a new channel when everything is ok,
				// or nil if it's necessary to quit.
				continue RecvLoop

			// receives and adds a chan to the list of notifyChans
			case notifyChan := <-acnl.subChan:
				if notifyChan != nil {
					acnl.l.WithFields(log.Fields{
						"tag": tag,
					}).Debug("new notify channel received")

					notifyChans = append(notifyChans, notifyChan)
				}

			case delivery, ok := <-deliveries:
				if !ok {
					acnl.l.WithFields(log.Fields{
						"tag": tag,
					}).Debug("delivery chan closed")

					continue RecvLoop
				}

				newNotifyChans, err := acnl.handleDelivery(delivery, notifyChans)
				if err != nil {
					acnl.l.WithError(err).WithFields(log.Fields{
						"tag": tag,
					}).Error("error during delivery handling")
				} else {
					notifyChans = newNotifyChans
				}
			}
		}
	}

	// closing all the notification channels
	for _, cnl := range notifyChans {
		close(cnl)
	}

	acnl.l.WithFields(log.Fields{
		"tag": "receiver-amqp",
	}).Infoln("AMQP receiver exiting")

	acnl.wg.Done()
}

// spawnReceiver spawns a goroutine which handles the reception of
// incoming messages from the NFVO on a dedicated queue.
// The receiver main channel is updated by setup() with a new
// consumer each time the connection is reestablished.
func (acnl *Channel) spawnReceiver() {
	acnl.wg.Add(1)

	go acnl.receiver()
}
