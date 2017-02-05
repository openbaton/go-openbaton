package amqp

import (
	"encoding/json"
	"time"

	"github.com/openbaton/go-openbaton/util"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func (acnl *Channel) connSetup() (*amqp.Connection, error) {
	tag := util.FuncName()

	acnl.l.WithFields(log.Fields{
		"tag": tag,
	}).Info("dialing AMQP")

	return amqp.DialConfig(acnl.cfg.connstr, acnl.cfg.cfg)	
}

func (acnl *Channel) getAMQPChan() (*amqp.Channel, error) {
	select {
	case acnl.chanReqChan <- struct{}{}:

	case <-time.After(time.Second):
		// the main loop is not listening
		return nil, ErrTimedOut
	}

	resp, ok := <-acnl.newChanChan
	if !ok {
		return nil, nil
	}

	return resp.Channel, resp.error
}

func (acnl *Channel) makeAMQPChan(conn *amqp.Connection) (*amqp.Channel, error) {
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

	return cnl, nil
}

func (acnl *Channel) register(conn *amqp.Connection) error {
	tag := util.FuncName()

	msg, err := json.Marshal(acnl.endpoint)
	if err != nil {
		return err
	}

	acnl.l.WithFields(log.Fields{
		"tag":      tag,
		"endpoint": string(msg),
	}).Info("sending a registering request to the NFVO")

	// use the newly instantiated connection to register
	cnl, err := acnl.makeAMQPChan(conn)
	if err != nil {
		return err
	}

	return acnl.publish(cnl, QueueVNFMRegister, msg)
}

func (acnl *Channel) setupQueues(cnl *amqp.Channel) error {
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
func (acnl *Channel) unregister(conn *amqp.Connection) (err error) {
	tag := util.FuncName()

	const Attempts = 2

	var msg []byte
	msg, err = json.Marshal(acnl.endpoint)
	if err != nil {
		return
	}

	acnl.l.WithFields(log.Fields{
		"tag":          tag,
		"max-attempts": Attempts,
		"endpoint":     string(msg),
	}).Debug("sending an unregistering request")

	for i := 0; i < Attempts; i++ {
		// Try to use the current connection the first time.
		// Recreate it otherwise

		if i > 0 {
			acnl.l.WithFields(log.Fields{
				"tag": tag,
				"try": i,
			}).Warn("attempting to re-initialize the connection")

			if _, err = acnl.connSetup(); err != nil {
				acnl.l.WithError(err).WithFields(log.Fields{
					"tag": tag,
					"try": i,
				}).Warn("conn setup failed")
				continue
			}
		}

		var cnl *amqp.Channel

		cnl, err = acnl.makeAMQPChan(conn)
		if err != nil {
			acnl.l.WithError(err).WithFields(log.Fields{
				"tag": tag,
				"try": i,
			}).Warn("chan setup failed")
			continue
		}

		defer cnl.Close()

		err = acnl.publish(cnl, QueueVNFMUnregister, msg)

		if err == nil {
			acnl.l.WithFields(log.Fields{
				"tag":     tag,
				"try":     i,
				"success": true,
			}).Info("endpoint unregister request successfully sent")
			break
		}

		acnl.l.WithFields(log.Fields{
			"tag":     tag,
			"try":     i,
			"success": false,
		}).Warn("endpoint unregister failed to send")
	}

	return
}
