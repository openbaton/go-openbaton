/*
 *  Copyright (c) 2017 Open Baton (http://openbaton.org)
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package plugin

import (
	"encoding/json"

	"github.com/openbaton/go-openbaton/util"
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
				req.DeliveryTag = delivery.DeliveryTag

				p.l.WithFields(log.Fields{
					"tag":        tag,
					"req-method": req.MethodName,
				}).Debug("received message")

				p.reqChan <- req

				p.l.WithFields(log.Fields{
					"tag":        tag,
					"req-method": req.MethodName,
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
