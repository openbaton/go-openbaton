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

package amqp

import (
	"errors"
	"time"

	"github.com/openbaton/go-openbaton/catalogue/messages"
	"github.com/openbaton/go-openbaton/vnfm/channel"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	ErrTimedOut = errors.New("timed out")
)

func (acnl *Channel) closeQueues() {
	acnl.setStatus(channel.Quitting)

	close(acnl.statusChan)
	close(acnl.sendQueue)
	close(acnl.subChan)
	close(acnl.quitChan)
	close(acnl.receiverDeliveryChan)

	// wait for all workers to quit
	acnl.wg.Wait()
}

func (acnl *Channel) setStatus(newStatus channel.Status) {
	for i := 0; i < acnl.numOfWorkers; i++ {
		acnl.statusChan <- newStatus
	}

	acnl.status = newStatus
}

// spawn spawns the main handler for AMQP communications.
func (acnl *Channel) spawn() error {
	errChan, err := acnl.setup()
	if err != nil {
		return err
	}

	acnl.register()

	acnl.spawnWorkers()
	acnl.spawnReceiver()
	acnl.setStatus(channel.Running)

	go func() {
		for {
			select {
			case <-acnl.quitChan:
				if err = acnl.unregister(); err != nil {
					acnl.l.WithFields(log.Fields{
						"tag": "channel-amqp",
						"err": err,
					}).Error("unregister failed")
				}

				if err = acnl.conn.Close(); err != nil {
					acnl.l.WithFields(log.Fields{
						"tag": "channel-amqp",
						"err": err,
					}).Error("closing Connection failed")

					acnl.closeQueues()
					return
				}

				acnl.l.WithFields(log.Fields{
					"tag": "channel-amqp",
				}).Info("initiating clean shutdown")

				// Close will cause the reception of nil on errChan.

			case amqpErr := <-errChan:
				// The connection closed cleanly after invoking Close().
				if amqpErr == nil {
					// notify the receiving end and listeners
					acnl.closeQueues()

					return
				}

				acnl.l.WithFields(log.Fields{
					"tag": "channel-amqp",
					"err": amqpErr,
				}).Error("received AMQP error for current connection")

				acnl.setStatus(channel.Reconnecting)

				// The connection crashed for some reason. Try to bring it up again.
				for {
					if errChan, err = acnl.setup(); err != nil {
						acnl.l.WithFields(log.Fields{
							"tag": "channel-amqp",
						}).Error("can't re-establish connection with AMQP; queues stalled. Retrying in 30 seconds.")
						time.Sleep(30 * time.Second)
					} else {
						acnl.setStatus(channel.Running)
						break
					}
				}

			}
		}
	}()

	return nil
}

// spawnReceiver spawns a goroutine which handles the reception of
// incoming messages from the NFVO on a dedicated queue.
// The receiver main channel is updated by setup() with a new
// consumer each time the connection is reestablished.
func (acnl *Channel) spawnReceiver() {
	acnl.wg.Add(1)

	go func() {
		acnl.l.WithFields(log.Fields{
			"tag": "receiver-amqp",
		}).Infoln("AMQP receiver starting")

		// list of channels to which incoming messages will be broadcasted.
		notifyChans := []chan<- messages.NFVMessage{}

		var deliveryChan <-chan amqp.Delivery
	RecvLoop:
		for {
			select {
			// setup delivers a new channel to this receiver, to
			// be listened for Deliveries.
			case deliveryChan = <-acnl.receiverDeliveryChan:
				if deliveryChan == nil {
					break RecvLoop
				}

				acnl.l.WithFields(log.Fields{
					"tag": "receiver-amqp",
				}).Debug("new delivery channel received")
				// chan updated

			// receives and adds a chan to the list of notifyChans
			case notifyChan := <-acnl.subChan:
				if notifyChan != nil {
					acnl.l.WithFields(log.Fields{
						"tag": "receiver-amqp",
					}).Debug("new notify channel received")

					notifyChans = append(notifyChans, notifyChan)
				}

			case delivery, ok := <-deliveryChan:
				if ok {
					msg, err := messages.Unmarshal(delivery.Body, messages.NFVO)
					if err != nil {
						acnl.l.WithFields(log.Fields{
							"tag": "receiver-amqp",
							"err": err,
						}).Error("message unmarshaling error")
						continue RecvLoop
					}

					acnl.l.WithFields(log.Fields{
						"tag": "receiver-amqp",
						"msg": msg,
					}).Debug("received message")

					last := 0
					for _, c := range notifyChans {
						select {
						// message sent successfully.
						case c <- msg:
							// keep the channel around for the next time
							notifyChans[last] = c
							last++

						// nobody is listening at the other end of the channel.
						case <-time.After(1 * time.Second):
							acnl.l.WithFields(log.Fields{
								"tag": "receiver-amqp",
							}).Debug("closing unresponsive notify channel")

							close(c)
						}
					}

					// notifyChans trimmed of dead chans
					notifyChans = notifyChans[:last]

					acnl.l.WithFields(log.Fields{
						"tag":          "receiver-amqp",
						"msg":          msg,
						"num-of-chans": last,
					}).Debug("message dispatched")

				} else {
					// make deliveryChan nil if someone closes it:
					// a closed channel always immediately returns a zero value, thus never
					// allowing the select to block.
					// A nil channel always blocks.

					acnl.l.WithFields(log.Fields{
						"tag": "receiver-amqp",
					}).Debug("delivery chan closed")

					deliveryChan = nil
				}
			}
		}

		// closing all the notification channels
		for _, cnl := range notifyChans {
			close(cnl)
		}

		acnl.l.WithFields(log.Fields{
			"tag": "receiver-amqp",
		}).Infoln("AMQP receiver exiting")

		acnl.wg.Done()
	}()
}
func (acnl *Channel) spawnWorkers() {
	acnl.wg.Add(acnl.numOfWorkers)
	for i := 0; i < acnl.numOfWorkers; i++ {
		go acnl.worker(i)
	}
}

func (acnl *Channel) worker(id int) {
	acnl.l.WithFields(log.Fields{
		"tag":       "worker-amqp",
		"worker-id": id,
	}).Debug("AMQP worker starting")

	status := channel.Stopped

	// explanation: a read on a nil channel will
	// block forever. This lambda ensures that we will accept jobs only
	// when the status is valid.
	work := func() chan *exchange {
		if status == channel.Running {
			return acnl.sendQueue
		}

		return nil
	}

WorkerLoop:
	for {
		select {
		// Updates the status. If it becomes Running, the next loop will accept incoming jobs again
		case status = <-acnl.statusChan:
			if status == channel.Quitting {
				break WorkerLoop
			}

		case exc := <-work():
			if exc.replyChan != nil { // RPC request
				resp, err := acnl.rpc(exc.queue, exc.msg)

				exc.replyChan <- response{resp, err}
			} else { //send only
				if err := acnl.publish(exc.queue, exc.msg); err != nil {
					acnl.l.WithFields(log.Fields{
						"tag":       "worker-amqp",
						"worker-id": id,
						"err":       err,
					}).Error("publish failed")
				}
			}
		}
	}

	acnl.l.WithFields(log.Fields{
		"tag":       "worker-amqp",
		"worker-id": id,
	}).Debug("AMQP worker stopping")

	acnl.wg.Done()
}
