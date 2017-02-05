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
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/openbaton/go-openbaton/util"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const pluginExchange = "plugin-exchange"

var (
	ErrNotInitialised = errors.New("not connected yet. Retry later")
	ErrTimeout        = errors.New("timed out")
)

// New creates a plugin from an implementation and plugin.Params.
// impl must be of a valid Plugin implementation type, like plugin.Driver.
func New(impl interface{}, p *Params) (Plugin, error) {
	tag := util.FuncName()

	if p.Workers < 1 {
		p.Workers = 10
	}

	plug := &plug{
		connstr:     util.AmqpUriBuilder(p.Username, p.Password, p.BrokerAddress, "", p.Port, false),
		params:      p,
		quitChan:    make(chan error),
		chanReqChan: make(chan struct{}, p.Workers+1),
		newChanChan: make(chan struct {
			*amqp.Channel
			error
		}),
	}

	if err := plug.initLogger(); err != nil {
		return nil, err
	}

	var rh reqHandler

	switch v := impl.(type) {
	case Driver:
		rh = driverHandler{v, plug.l}

	// in case we are reinitialising the plugin
	case reqHandler:
		rh = v

	default:
		plug.l.WithField("tag", tag).Panicf("invalid plugin implementation %T", impl)
	}

	plug.rh = rh

	return plug, nil
}

type plug struct {
	connstr string

	l        *log.Logger
	e        logData
	params   *Params
	quitChan chan error

	// because the connection is private,
	// to get a channel is necessary to request it
	// through this channel.
	// The request will be received, and the channel will be sent throgh the newChanChan channel
	chanReqChan chan struct{}
	newChanChan chan struct {
		*amqp.Channel
		error
	}

	rh      reqHandler
	stopped bool
	wg      sync.WaitGroup
}

func (p *plug) ChannelAccessor() func() (*amqp.Channel, error) {
	return p.getAMQPChan
}

func (p *plug) Logger() *log.Logger {
	return p.l
}

func (p *plug) Serve() error {
	tag := util.FuncName()

	// reinit the plugin if already stopped
	if p.stopped {
		panic("plugin already stopped")
	}

	p.l.WithFields(log.Fields{
		"tag":    tag,
		"params": *p.params,
	}).Debug("plugin starting")

	p.spawnWorkers()

	exiting := false

MainLoop:
	for {
		conn, err := p.connSetup()
		if err != nil {
			return err
		}

		errChan := makeErrChan(conn)

		for {
			select {
			case <-p.quitChan:
				if err = conn.Close(); err != nil {
					p.l.WithError(err).WithFields(log.Fields{
						"tag": tag,
					}).Error("closing Connection failed")

					p.closeQueues()

					// send the error to stop
					p.quitChan <- err
					return nil
				}

				p.l.WithFields(log.Fields{
					"tag": tag,
				}).Info("initiating clean shutdown")

				exiting = true

				// Close will cause the reception of nil on errChan.

			// some worker wants a Channel
			case <-p.chanReqChan:
				var cnl *amqp.Channel
				var err error

				// avoid trying to create a channel when exiting
				if !exiting {
					cnl, err = p.makeAMQPChan(conn)
				}

				p.newChanChan <- struct {
					*amqp.Channel
					error
				}{cnl, err}

				// after sending the response, check if it was ok.
				// If there was an error, the client has issues
				if err != nil {
					// the connection is broken.
					// create a new one.
					errChan = nil // avoid receiving anything
					continue MainLoop
				}

			case amqpErr := <-errChan:
				// The connection closed cleanly after invoking Close().
				if amqpErr == nil {
					// notify the receiver and workers
					p.closeQueues()

					p.wg.Wait()

					// send nil to Stop
					close(p.quitChan)

					p.l.WithFields(log.Fields{
						"tag": tag,
					}).Debug("main loop quitting")

					break MainLoop
				}

				p.l.WithError(amqpErr).WithFields(log.Fields{
					"tag": tag,
				}).Error("received AMQP error for current connection")

				// The connection crashed for some reason. Try to bring it up again.
				for {
					conn, err = p.connSetup()
					if err != nil {
						p.l.WithError(err).WithFields(log.Fields{
							"tag": tag,
						}).Error("can't re-establish connection with AMQP; queues stalled. Retrying in 30 seconds.")

						time.Sleep(30 * time.Second)
					} else {
						errChan = makeErrChan(conn)
					}
				}

			}
		}
	}

	return nil
}

func (p *plug) Stop() error {
	tag := util.FuncName()

	defer p.deinitLogger()

	if p.stopped {
		return fmt.Errorf("plugin %s already stopped", p.params.Name)
	}

	// first step: signal the main routine to quit.
	select {
	case p.quitChan <- nil:

	case <-time.After(time.Second):
		return errors.New("the plugin is not listening")
	}

	// second step: wait for it to quit
	select {
	case err := <-p.quitChan:
		if err != nil {
			return err
		}
	case <-time.After(1 * time.Minute):
		return errors.New("the plugin refused to quit")
	}

	p.stopped = true

	p.l.WithFields(log.Fields{
		"tag": tag,
	}).Info("plugin stopped cleanly")

	return nil
}

func (p *plug) Type() string {
	return p.rh.Type()
}

func (p *plug) closeQueues() {
	// closes the new chan channel
	close(p.newChanChan)

	p.wg.Wait()
}

func (p *plug) id() string {
	return fmt.Sprintf("%s.%s.%s", p.rh.QueueTag(), p.params.Type, p.params.Name)
}

func makeErrChan(conn *amqp.Connection) chan *amqp.Error {
	return conn.NotifyClose(make(chan *amqp.Error))
}

type reqHandler interface {
	Handle(call string, args []json.RawMessage) (interface{}, error)
	QueueTag() string
	Type() string
}
