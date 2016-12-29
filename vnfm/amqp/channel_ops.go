package amqp

import (
	"errors"
	"time"

	"github.com/mcilloni/go-openbaton/catalogue"
	"github.com/mcilloni/go-openbaton/catalogue/messages"
	"github.com/mcilloni/go-openbaton/vnfm/channel"
	"github.com/streadway/amqp"
)

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

func (acnl *amqpChannel) publish(queue string, msg []byte) error {
	return acnl.cnl.Publish(
		acnl.cfg.exchange.name, 
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: msg,
		},
	)
}

func (acnl *amqpChannel) rpc(queue string, msg []byte) ([]byte, error) {
	replyQueue, err := acnl.temporaryQueue()
	if err != nil {
		return nil, err
	}

	deliveries, err := acnl.cnl.Consume(
		replyQueue, // queue
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

	corrID := string(catalogue.GenerateID())

	err = acnl.cnl.Publish(
		acnl.cfg.exchange.name, 
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			CorrelationId: corrID,
			ReplyTo: replyQueue,
			Body: msg,
		},
	)
	if err != nil {
		return nil, err
	}

	for delivery := range deliveries {
		if delivery.CorrelationId == corrID {
			return delivery.Body, nil
		}
	}

	return nil, errors.New("no reply received")
}

func (acnl *amqpChannel) temporaryQueue() (string, error) {
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