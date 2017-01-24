package amqp

import (
	"errors"
	"time"

	"github.com/mcilloni/go-openbaton/catalogue/messages"
	"github.com/mcilloni/go-openbaton/util"
	"github.com/mcilloni/go-openbaton/vnfm/channel"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Close closes the Channel.
func (acnl *Channel) Close() error {
	acnl.quitChan <- struct{}{}

	select {
	case <-acnl.quitChan:
		return nil

	case <-time.After(1 * time.Minute):
		return errors.New("timed out afer waiting for AMQP handler loop to close")
	}
}

// AMQPExchange returns a string containing the default exchange.
func (acnl *Channel) AMQPExchange() string {
	return acnl.cfg.exchange.name
} 

// Exchange executes an RPC call to a given queue on the default exchange.
func (acnl *Channel) Exchange(queue string, msg []byte) ([]byte, error) {
	respChan := make(chan response)

	acnl.sendQueue <- &exchange{queue, msg, respChan}

	resp := <-respChan
	return resp.msg, resp.error
}

// NFVOExchange delivers a message to the NFVO through an RPC call, and awaits a response.
func (acnl *Channel) NFVOExchange(msg messages.NFVMessage) (messages.NFVMessage, error) {
	msgBytes, err := messages.Marshal(msg)
	if err != nil {
		return nil, err
	}

	retBytes, err := acnl.Exchange(QueueVNFMCoreActionsReply, msgBytes)
	if err != nil {
		return nil, err
	}

	return messages.Unmarshal(retBytes, messages.NFVO)
}

// NFVOSend delivers a message to the NFVO.
func (acnl *Channel) NFVOSend(msg messages.NFVMessage) error {
	msgBytes, err := messages.Marshal(msg)
	if err != nil {
		return err
	}

	return acnl.Send(QueueVNFMCoreActions, msgBytes)
}

// NotifyReceived creates and returns a channel of NFVMessage; every received message
// will be broadcasted to every channel created by this function.
// If nobody is listening on the receiving channel, the channel will be dropped.
func (acnl *Channel) NotifyReceived() (<-chan messages.NFVMessage, error) {
	notifyChan := make(chan messages.NFVMessage, 5)

	acnl.subChan <- notifyChan

	return notifyChan, nil
}

// Send sends a message to a given queue on the default exchange.
func (acnl *Channel) Send(queue string, msg []byte) error {
	acnl.sendQueue <- &exchange{queue, msg, nil}

	return nil
}

// Status returns the current Status of the Channel.
func (acnl *Channel) Status() channel.Status {
	return acnl.status
}

func (acnl *Channel) publish(queue string, msg []byte) error {
	return acnl.cnl.Publish(
		acnl.cfg.exchange.name,
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)
}

func (acnl *Channel) rpc(queue string, msg []byte) ([]byte, error) {
	replyQueue, err := acnl.temporaryQueue()
	if err != nil {
		return nil, err
	}

	deliveries, err := acnl.cnl.Consume(
		replyQueue, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, err
	}

	corrID := util.GenerateID()

	acnl.l.WithFields(log.Fields{
		"tag":            "channel-amqp-rpc",
		"corr-id":        corrID,
		"reply-to-queue": replyQueue,
	}).Debug("sending RPC publishing")

	err = acnl.cnl.Publish(
		acnl.cfg.exchange.name,
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrID,
			ReplyTo:       replyQueue,
			Body:          msg,
		},
	)
	if err != nil {
		return nil, err
	}

	timeout := time.After(DefaultTimeout)

DeliveryLoop:
	for {
		select {
		case <-timeout:
			break DeliveryLoop
			
		case delivery, ok := <-deliveries:
			if !ok {
				break DeliveryLoop
			}
			
			acnl.l.WithFields(log.Fields{
				"tag": "channel-amqp-rpc",
			}).Debug("received delivery")

			if delivery.CorrelationId == corrID {
				return delivery.Body, nil
			}
		}
	}

	return nil, errors.New("no reply received")
}

func (acnl *Channel) temporaryQueue() (string, error) {
	queue, err := acnl.cnl.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return "", err
	}

	return queue.Name, nil
}
