package amqp

import (
    "time"

	"github.com/mcilloni/go-openbaton/catalogue/messages"
)

type amqpChannel struct {}

func newChannel() (*amqpChannel, error) {
    return &amqpChannel{}, nil
}

func (cnl *amqpChannel) Close() error {
    return nil
}

func (cnl *amqpChannel) Exchange(msg messages.NFVMessage, timeout time.Duration) (messages.NFVMessage, error) {
    return nil, nil
}
	
func (cnl *amqpChannel) ExchangeStrings(msg, queue string, timeout time.Duration) (string, error) {
    return "", nil
}

func (cnl *amqpChannel) NotifyReceived() (<-chan messages.NFVMessage, error) {
    return nil, nil
}

func (cnl *amqpChannel) Send(msg messages.NFVMessage) error {
    return nil
}