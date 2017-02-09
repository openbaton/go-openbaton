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

	"github.com/openbaton/go-openbaton/util"
	"github.com/openbaton/go-openbaton/vnfm/channel"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	ErrTimedOut = errors.New("timed out")
)

func (acnl *Channel) closeQueues() {
	close(acnl.newChanChan)
	close(acnl.statusChan)
	close(acnl.sendQueue)
	close(acnl.subChan)
	close(acnl.quitChan)

	// wait for all workers to quit
	acnl.wg.Wait()
}

// mainLoop handles the requests.
func (acnl *Channel) mainLoop(conn *amqp.Connection) {
	tag := util.FuncName()

	errChan := makeErrChan(conn)

	first := true

MainLoop:
	for {
		var err error

		// If this is the first time, then use the connection provided by the caller
		// It this is not the first time, recreate a new connection
		if !first {
			conn, err = acnl.connSetup()
			if err != nil {
				acnl.l.WithError(err).WithFields(log.Fields{
					"tag": tag,
				}).Error("can't re-establish connection with AMQP; queues stalled. Retrying in 30 seconds.")
				time.Sleep(30 * time.Second)
				continue MainLoop
			} else {
				acnl.l.WithFields(log.Fields{
					"tag": tag,
				}).Debug("AMQP connection reestablished")

				// update the errChan with a new one
				errChan = makeErrChan(conn)

				// resume normal operations
			}
		} else {
			first = false
		}

		acnl.setStatus(channel.Running)

		for {
			select {
			case <-acnl.quitChan:
				acnl.setStatus(channel.Quitting)

				if err := acnl.unregister(conn); err != nil {
					acnl.l.WithError(err).WithFields(log.Fields{
						"tag": tag,
					}).Error("unregister failed")
				}

				if err := conn.Close(); err != nil {
					acnl.l.WithError(err).WithFields(log.Fields{
						"tag": tag,
					}).Error("closing Connection failed")

					acnl.closeQueues()
					return
				}

				acnl.l.WithFields(log.Fields{
					"tag": tag,
				}).Info("initiating clean shutdown")

				// Close will cause the reception of nil on errChan.

			case <-acnl.chanReqChan:
				// somebody wants a new channel.

				var resp struct {
					*amqp.Channel
					error
				}

				// if we are quitting, send nil back
				if acnl.status != channel.Quitting {
					resp.Channel, resp.error = acnl.makeAMQPChan(conn)
				}

				acnl.newChanChan <- resp

				// after sending the response, check if it was ok.
				// If there was an error, the connection has issues
				if err != nil {
					// the connection is broken.
					// create a new one.
					errChan = nil // avoid receiving anything
					continue MainLoop
				}

			case amqpErr := <-errChan:
				// The connection closed cleanly after invoking Close().
				if amqpErr == nil {
					acnl.l.WithFields(log.Fields{
						"tag": tag,
					}).Debug("shutting down workers...")

					// notify the receiving end and listeners
					acnl.closeQueues()

					return
				}

				acnl.l.WithError(amqpErr).WithFields(log.Fields{
					"tag": tag,
				}).Error("received AMQP error for current connection")

				acnl.setStatus(channel.Reconnecting)

				// The connection crashed for some reason. Try to bring it up again.

				continue MainLoop
			}
		}
	}

}

func (acnl *Channel) setStatus(newStatus channel.Status) {
	for i := 0; i < acnl.numOfWorkers; i++ {
		acnl.statusChan <- newStatus
	}

	acnl.status = newStatus
}

// spawn spawns the main handler for AMQP communications.
func (acnl *Channel) spawn() error {
	//tag := util.FuncName()

	// I'm allocating a new connection here to return an error in case
	// the parameters are incorrect, instead of spawning a routine
	conn, err := acnl.connSetup()
	if err != nil {
		return err
	}

	if err := acnl.register(conn); err != nil {
		return err
	}

	acnl.spawnWorkers()
	acnl.spawnReceiver()

	// give the main loop the newly allocated conn
	go acnl.mainLoop(conn)

	return nil
}

func makeErrChan(conn *amqp.Connection) chan *amqp.Error {
	return conn.NotifyClose(make(chan *amqp.Error))
}
