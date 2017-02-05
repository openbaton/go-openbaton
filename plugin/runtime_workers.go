package plugin

import (
	"encoding/json"

	"github.com/openbaton/go-openbaton/util"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func (p *plug) spawnWorkers() {
	tag := util.FuncName()

	p.l.WithFields(log.Fields{
		"tag":            tag,
		"num-of-workers": p.params.Workers,
	}).Debug("spawning workers")

	p.wg.Add(p.params.Workers)
	for i := 0; i < p.params.Workers; i++ {
		go p.worker(i)
	}
}

func (p *plug) temporaryQueue() (string, error) {
	queue, err := p.cnl.QueueDeclare(
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

func (p *plug) worker(id int) {
	tag := util.FuncName()

	p.l.WithFields(log.Fields{
		"tag":       tag,
		"worker-id": id,
	}).Debug("worker is starting")

	for req := range p.reqChan {
		result, err := p.rh.Handle(req.MethodName, req.Parameters)

		var resp response
		if err != nil {
			// The NFVO expects a Java Exception;
			// This type switch checks if the error is not one of the special
			// Java-compatible types already and wraps it.
			switch err.(type) {
			case plugError:
				resp.Exception = err

			case DriverError:
				resp.Exception = err

			// if the error is not a special plugin error, than wrap it:
			// the nfvo expects a Java exception.
			default:
				resp.Exception = plugError{err.Error()}
			}
		} else {
			resp.Answer = result
		}

		bResp, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			p.l.WithError(err).WithFields(log.Fields{
				"tag":       tag,
				"worker-id": id,
			}).Error("failure while serialising response")
			continue
		}

		err = p.cnl.Publish(
			pluginExchange,
			req.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: req.CorrID,
				Body:          bResp,
			},
		)

		if err != nil {
			p.l.WithError(err).WithFields(log.Fields{
				"tag":          tag,
				"worker-id":    id,
				"reply-queue":  req.ReplyTo,
				"reply-corrid": req.CorrID,
			}).Error("failure while replying")
			continue
		}

		p.l.WithError(resp.Exception).WithFields(log.Fields{
			"tag":          tag,
			"worker-id":    id,
			"reply-queue":  req.ReplyTo,
			"reply-corrid": req.CorrID,
		}).Info("response sent")

		// IMPORTANT: Acknowledge the received delivery!
		// The VimDriverCaller executor thread of the NFVO
		// will perpetually sleep when trying to publish the
		// next request if this step is omitted.
		if err := p.cnl.Ack(req.DeliveryTag, false); err != nil {
			p.l.WithError(err).WithFields(log.Fields{
				"tag":                tag,
				"worker-id":          id,
				"reply-queue":        req.ReplyTo,
				"reply-corrid":       req.CorrID,
				"reply-delivery_tag": req.DeliveryTag,
			}).Error("failure while acknowledging the last delivery")
			continue
		}
	}

	p.l.WithFields(log.Fields{
		"tag":       tag,
		"worker-id": id,
	}).Debug("worker is stopping")

	p.wg.Done()
}
