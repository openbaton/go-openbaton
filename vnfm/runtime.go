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

package vnfm

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/openbaton/go-openbaton/catalogue/messages"
	"github.com/openbaton/go-openbaton/vnfm/channel"
	"github.com/openbaton/go-openbaton/vnfm/config"
	log "github.com/sirupsen/logrus"
)

var impls = make(map[string]channel.Driver)

// Register registers a channel.Driver. Invoke this in an init() method of a driver package.
func Register(name string, driver channel.Driver) {
	if _, ok := impls[name]; ok {
		panic(fmt.Sprintf("trying to register driver of type %T with already existing name '%s'", driver, name))
	}

	if driver == nil {
		panic("nil driver")
	}

	impls[name] = driver
}

// VNFM represents a VNFM instance.
type VNFM interface {
	// ChannelAccessor returns a function that returns the underlying channel.Channel of this VNFM.
	// If Serve() hasn't been called, the returned function will lock until it's ready.
	ChannelAccessor() func() (channel.Channel, error)

	// Logger returns a logrus logger instance.
	Logger() *log.Logger

	// Serve launches the VNFM. No error will be returned after a valid initialization. See the logfile or the return value of Stop().
	Serve() error

	// Stop signals the VNFM to quit. It returns when the VNFM quits or after a timeout.
	Stop() error
}

// New returns a new VNFM. implName must be a string representing a Driver previously registered by Register().
func New(implName string, handler Handler, config *config.Config) (VNFM, error) {
	if _, ok := impls[implName]; !ok {
		return nil, fmt.Errorf("no implementation available for %s. Have you forgot to import its package?", implName)
	}

	logger := log.New()

	logger.Formatter = &log.TextFormatter{
		ForceColors:   config.LogColors,
		DisableColors: !config.LogColors,
	}

	logger.Level = config.LogLevel

	if config.LogFile != "" {
		file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
		if err != nil {
			return nil, fmt.Errorf("couldn't open the log file %s: %s", config.LogFile, err.Error())
		}

		logger.Out = file
	} else {
		logger.Out = terminalWriter()
	}

	return &vnfm{
		cnlCond:  sync.NewCond(&sync.Mutex{}),
		hnd:      handler,
		implName: implName,
		conf:     config,
		l:        logger,
		quitChan: make(chan error, 1), // do not block on send
	}, nil
}

type vnfm struct {
	cnl      channel.Channel
	cnlErr   error
	cnlCond  *sync.Cond
	conf     *config.Config
	hnd      Handler
	implName string
	l        *log.Logger
	msgChan  <-chan messages.NFVMessage
	quitChan chan error
	wg       sync.WaitGroup
}

func (vnfm *vnfm) ChannelAccessor() func() (channel.Channel, error) {
	return vnfm.channel
}

func (vnfm *vnfm) Logger() *log.Logger {
	return vnfm.l
}

func (vnfm *vnfm) Serve() (err error) {
	vnfm.cnl, err = impls[vnfm.implName].Init(vnfm.conf, vnfm.l)

	if err != nil {
		vnfm.cnl = nil
		vnfm.cnlErr = errors.New("fetching channel failed")
	}

	// tell the Channel() listeners that we have a result,
	// either good or bad
	vnfm.cnlCond.Broadcast()

	if err != nil {
		return
	}

	defer func() {
		r := recover()

		// Check if if a file has been opened in New.
		if file, ok := vnfm.l.Out.(*os.File); ok {
			if err := file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error while closing logfile: %v\n", err)
			}
		}

		// answering the channel signals Stop() that we're quitting
		err := vnfm.cnl.Close()

		if err == nil {
			// if the channel closed politely, the workers will be quitting by now;
			// otherwise, they will be killed when main exits.
			vnfm.wg.Wait()
		}

		vnfm.quitChan <- err

		if r != nil {
			vnfm.l.WithFields(log.Fields{
				"tag":         "vnfm-serve-on_exit",
				"stack-trace": string(debug.Stack()),
			}).Panic(r)
		}
	}()

	if vnfm.msgChan, err = vnfm.cnl.NotifyReceived(); err != nil {
		return
	}

	vnfm.spawnWorkers()

	// wait for Stop()
	<-vnfm.quitChan

	return
}

func (vnfm *vnfm) SetLogger(log *log.Logger) {
	vnfm.l = log
}

func (vnfm *vnfm) Stop() error {
	select {
	case vnfm.quitChan <- nil:

	case <-time.After(time.Second):
		return errors.New("the VNFM is not listening")
	}

	select {
	case err := <-vnfm.quitChan:
		return err
	case <-time.After(1 * time.Minute):
		return errors.New("the VNFM refused to quit")
	}
}

func (vnfm *vnfm) channel() (channel.Channel, error) {
	// Use a condition: get the lock, check if there is a channel or an error, and then Wait for it.
	vnfm.cnlCond.L.Lock()
	defer vnfm.cnlCond.L.Unlock()

	for vnfm.cnl == nil && vnfm.cnlErr == nil {
		vnfm.cnlCond.Wait()
	}

	return vnfm.cnl, vnfm.cnlErr
}

func (vnfm *vnfm) spawnWorkers() {
	const NumWorkers = 5

	vnfm.wg.Add(NumWorkers)

	for i := 0; i < NumWorkers; i++ {
		go (&worker{vnfm, i}).work()
	}
}
