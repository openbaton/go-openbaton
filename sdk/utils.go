package sdk

import (
	"github.com/op/go-logging"
	"os"
	"strings"
	"math/rand"
	"runtime/debug"
	"encoding/json"
	"github.com/streadway/amqp"
	"github.com/openbaton/go-openbaton/catalogue"
	"github.com/pkg/errors"
	"fmt"
)

var log *logging.Logger

type SdkError struct {
	err error
}

func (e *SdkError) Error() string {
	return fmt.Sprintf("SDK-ERROR: %s", e.err)
}

func NewSdkError(msg string) *SdkError {
	return &SdkError{errors.New(msg)}
}

//Obtain the Logger pre-formatted
func GetLogger(name string, levelStr string) (*logging.Logger) {
	if log != nil {
		return log
	}
	level := toLogLevel(levelStr)
	log = logging.MustGetLogger(name)
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05} [%{level:.4s}] %{module:6.10s} -> %{longfunc:10.10s} ▶ %{color:reset} %{message}`,
	)
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetLevel(level, "")
	logging.SetBackend(backendFormatter)

	return log
}

func toLogLevel(lvlStr string) (lvl logging.Level) {
	switch strings.ToUpper(lvlStr) {
	case "DEBUG":
		lvl = logging.DEBUG

	case "INFO":
		lvl = logging.INFO

	case "WARN":
		lvl = logging.WARNING

	case "ERROR":
		lvl = logging.ERROR

	case "FATAL":
	case "CRITICAL":
	case "PANIC":
		lvl = logging.CRITICAL

	default:
		lvl = logging.DEBUG
	}

	return
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// Execute a AMQP RPC call to a specific queue
func Rpc(queue string, message interface{}, conn *amqp.Connection, l *logging.Logger) ([]byte, error) {

	l.Info("Executing RPC to queue: %s", queue)
	l.Debug("Getting Channel for RPC")
	channel, err := conn.Channel()
	defer channel.Close()
	l.Debug("Got Channel for RPC")

	var q amqp.Queue
	var msgs <-chan amqp.Delivery

	l.Debug("Declaring Queue for RPC")
	q, err = channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		debug.PrintStack()
		l.Errorf("Failed to declare a queue: %v", err)
		return nil, err
	}
	l.Debug("Declared Queue for RPC")
	l.Debug("Registering consumer for RPC")
	msgs, err = channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		debug.PrintStack()
		l.Errorf("Failed to register a consumer: %v", err)
		return nil, err
	}

	l.Debug("Registered consumer for RPC")
	corrId := randomString(32)

	mrs, err := json.Marshal(message)
	if err != nil {
		l.Errorf("Error while marshaling: %v", err)
		return nil, err
	}
	l.Debug("Publishing message to queue %s", queue)
	err = channel.Publish(
		OpenbatonExchangeName, // exchange
		queue,                 // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType:   AmqpContentType,
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(mrs),
		})

	if err != nil {
		l.Errorf("Failed to publish a message")
		return nil, err
	}
	l.Debugf("Published message to queue %s", queue)

	for d := range msgs {
		if corrId == d.CorrelationId {
			l.Debug("Received Response")
			return d.Body, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Not found message with correlationId [%s]", corrId))
}

// Send message to a specific queue
func SendMsg(queue string, message []byte, channel *amqp.Channel, logger *logging.Logger) (error) {
	err := channel.Publish(
		OpenbatonExchangeName,
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: AmqpContentType,
			Body:        []byte(message),
		})

	if err != nil {
		logger.Errorf("Failed to publish a message: %v", err)
		return err
	}
	return nil
}

// Vim Driver Error
type DriverError struct {
	Message string    `json:"detailMessage"`
	*catalogue.Server `json:"server"`
}

// Error returns a description of the error.
func (e DriverError) Error() string {
	return e.Message + " on server " + e.Server.Name
}
