package plugin

import (
	"encoding/json"

	"github.com/mcilloni/go-openbaton/util"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func (p *plug) receiver() {
	tag := util.FuncName()

	p.l.WithFields(log.Fields{
		"tag": tag,
	}).Infoln("AMQP receiver starting")

	// channel from which receive Deliveries
	var deliveryChan <-chan amqp.Delivery

RecvLoop:
	for {
		select {
		// setup delivers a new channel to this receiver, to
		// be listened for Deliveries.
		case deliveryChan = <-p.receiverDeliveryChan:
			if deliveryChan == nil {
				// quitting
				break RecvLoop
			}

			p.l.WithFields(log.Fields{
				"tag": tag,
			}).Debug("new delivery channel received")
			// chan updated

		case delivery, ok := <-deliveryChan:
			if ok {
				var req request
				if err := json.Unmarshal(delivery.Body, &req); err != nil {
					p.l.WithError(err).WithFields(log.Fields{
						"tag": tag,
					}).Error("message unmarshaling error")
					continue RecvLoop
				}

				req.ReplyTo = delivery.ReplyTo
				req.CorrID = delivery.CorrelationId

				p.l.WithFields(log.Fields{
					"tag": tag,
				}).Debug("received message")

				p.reqChan <- req

				p.l.WithFields(log.Fields{
					"tag": tag,
					"req": req,
				}).Debug("request dispatched")

			} else {
				// make deliveryChan nil if someone closes it:
				// a closed channel always immediately returns a zero value, thus never
				// allowing the select to block (infinite loop).
				// A nil channel always blocks.

				p.l.WithFields(log.Fields{
					"tag": tag,
				}).Debug("delivery chan closed")

				deliveryChan = nil
			}
		}
	}

	p.l.WithFields(log.Fields{
		"tag": tag,
	}).Infoln("AMQP receiver exiting")

	p.wg.Done()
}

func (p *plug) spawnReceiver() {
	p.wg.Add(1)

	go p.receiver()
}
