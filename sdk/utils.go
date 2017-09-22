package sdk

import (
	logging "github.com/op/go-logging"
	"os"
	"strings"
	"math/rand"
	"encoding/json"
	"github.com/streadway/amqp"
	"github.com/openbaton/go-openbaton/catalogue"
)

var log *logging.Logger

func GetLogger(name string, levelStr string) (*logging.Logger) {
	if log != nil {
		return log
	}
	level := toLogLevel(levelStr)
	log = logging.MustGetLogger(name)
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
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


func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func ExecuteRpc(queue string, message interface{}, channel *amqp.Channel, l *logging.Logger) (<- chan amqp.Delivery, string, error) {
	q, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)

	if err != nil{
		l.Errorf("Failed to declare a queue")
		return nil, "", err
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
	if err != nil{
		l.Errorf("Failed to register a consumer")
		return nil, "", err
	}


	corrId := RandomString(32)

	mrs, err := json.Marshal(message)
	if err != nil {
		l.Errorf("Error while marshaling: %v", err)
		return nil, "", err
	}
	err = channel.Publish(
		"openbaton-exchange",          // exchange
		queue, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(mrs),
		})

	if err != nil{
		l.Errorf("Failed to publish a message")
		return nil, "", err
	}
	return msgs, corrId, nil
}

type DriverError struct {
	Message           string `json:"detailMessage"`
	*catalogue.Server `json:"server"`
}

// Error returns a description of the error.
func (e DriverError) Error() string {
	return e.Message + " on server " + e.Server.Name
}
