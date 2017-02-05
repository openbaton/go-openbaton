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

func (p *plug) setupDeliveries(cnl *amqp.Channel) (<-chan amqp.Delivery, error) {
	// setup incoming deliveries
	return cnl.Consume(
		p.id(), // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}

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

func (p *plug) worker(id int) {
	tag := util.FuncName()

	p.l.WithFields(log.Fields{
		"tag":       tag,
		"worker-id": id,
	}).Debug("worker is starting")

WorkerLoop:
	for {
		cnl, err := p.getAMQPChan()
		if err != nil {
			p.l.WithError(err).WithFields(log.Fields{
				"tag":       tag,
				"worker-id": id,
			}).Error("failure while getting an AMQP channel")

			// retry
			continue WorkerLoop
		}

		// no more channels
		if cnl == nil {
			break WorkerLoop
		}

		errChan := cnl.NotifyClose(make(chan *amqp.Error))

		p.l.WithFields(log.Fields{
			"tag":       tag,
			"worker-id": id,
		}).Debug("received an AMQP channel")

		deliveries, err := p.setupDeliveries(cnl)
		if err != nil {
			p.l.WithError(err).WithFields(log.Fields{
				"tag":       tag,
				"worker-id": id,
			}).Error("error while setting up consumer")

			// get a new channel
			continue WorkerLoop
		}

		p.l.WithFields(log.Fields{
			"tag":       tag,
			"worker-id": id,
		}).Debug("set up consumer")

	ServeLoop:
		for {
			select {
			case err := <-errChan:
				if err != nil {
					p.l.WithError(err).WithFields(log.Fields{
						"tag":       tag,
						"worker-id": id,
					}).Error("channel failure")
				}

				// the channel broke.
				continue WorkerLoop

			case delivery, ok := <-deliveries:
				if !ok {
					break WorkerLoop
				}

				p.l.WithFields(log.Fields{
					"tag":          tag,
					"worker-id":    id,
					"reply-queue":  delivery.ReplyTo,
					"reply-corrid": delivery.CorrelationId,
				}).Debug("received delivery")

				var req request
				if err := json.Unmarshal(delivery.Body, &req); err != nil {
					p.l.WithError(err).WithFields(log.Fields{
						"tag":       tag,
						"worker-id": id,
					}).Error("message unmarshaling error")
					continue ServeLoop
				}

				p.l.WithFields(log.Fields{
					"tag":             tag,
					"worker-id":       id,
					"req-method_name": req.MethodName,
				}).Info("deserialised request")

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
					continue ServeLoop
				}

				err = cnl.Publish(
					pluginExchange,
					delivery.ReplyTo,
					false,
					false,
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: delivery.CorrelationId,
						Body:          bResp,
					},
				)

				if err != nil {
					p.l.WithError(err).WithFields(log.Fields{
						"tag":          tag,
						"worker-id":    id,
						"reply-queue":  delivery.ReplyTo,
						"reply-corrid": delivery.CorrelationId,
					}).Error("failure while replying")
					continue ServeLoop
				}

				p.l.WithError(resp.Exception).WithFields(log.Fields{
					"tag":          tag,
					"worker-id":    id,
					"reply-queue":  delivery.ReplyTo,
					"reply-corrid": delivery.CorrelationId,
				}).Info("response sent")

				// IMPORTANT: Acknowledge the received delivery!
				// The VimDriverCaller executor thread of the NFVO
				// will perpetually sleep when trying to publish the
				// next request if this step is omitted.
				if err := cnl.Ack(delivery.DeliveryTag, false); err != nil {
					p.l.WithError(err).WithFields(log.Fields{
						"tag":                tag,
						"worker-id":          id,
						"reply-queue":        delivery.ReplyTo,
						"reply-corrid":       delivery.CorrelationId,
						"reply-delivery_tag": delivery.DeliveryTag,
					}).Error("failure while acknowledging the last delivery")
					continue ServeLoop
				}
			}
		}
	}

	p.l.WithFields(log.Fields{
		"tag":       tag,
		"worker-id": id,
	}).Debug("worker is stopping")

	p.wg.Done()
}
