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
	"fmt"
	"sync"
	"time"

	"github.com/openbaton/go-openbaton/catalogue"
	"github.com/openbaton/go-openbaton/catalogue/messages"
	"github.com/openbaton/go-openbaton/util"
	"github.com/openbaton/go-openbaton/vnfm/channel"
	"github.com/openbaton/go-openbaton/vnfm/config"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type exchange struct {
	queue     string
	msg       []byte
	replyChan chan response
}

type response struct {
	msg []byte
	error
}

// A Channel is a control structure to handle an AMQP connection.
// The main logic is handled in an event loop, which is fed using Go channels through
// the amqpChannel methods.
type Channel struct {
	cfg struct {
		connstr string
		cfg     amqp.Config

		exchange struct {
			name    string
			durable bool
		}

		queues struct {
			autodelete, exclusive bool

			generic string
		}

		vnfmType, vnfmEndpoint, vnfmDescr string
	}

	endpoint *catalogue.Endpoint

	l            *log.Logger
	numOfWorkers int

	// because the connection is private,
	// to get a channel is necessary to request it
	// through this channel.
	// The request will be received, and the channel will be sent throgh the newChanChan channel
	chanReqChan chan struct{}
	newChanChan chan struct {
		*amqp.Channel
		error
	}

	quitChan   chan struct{}
	sendQueue  chan *exchange
	status     channel.Status
	statusChan chan channel.Status
	subChan    chan chan messages.NFVMessage

	wg sync.WaitGroup
}

func newChannel(config *config.Config, l *log.Logger) (*Channel, error) {
	props := config.Properties

	acnl := &Channel{
		l: l,

		quitChan: make(chan struct{}),
		status:   channel.Stopped,
		subChan:  make(chan chan messages.NFVMessage),
	}

	acnl.cfg.vnfmDescr = config.Description
	acnl.cfg.vnfmEndpoint = config.Endpoint
	acnl.cfg.vnfmType = config.Type
	acnl.cfg.queues.generic = fmt.Sprintf("nfvo.%s.actions", config.Type)

	// defaults
	host := "localhost"
	port := 5672
	username := ""
	password := ""
	vhost := ""
	heartbeat := 60
	exchangeName := ExchangeDefault
	exchangeDurable := true
	queuesExclusive := false
	queuesAutodelete := true

	workers, jobQueueSize := 5, 20

	if sect, ok := props.Section("amqp"); ok {
		acnl.l.WithFields(log.Fields{
			"tag": "channel-amqp-config",
		}).Info("found AMQP section in config")

		host, _ = sect.ValueString("host", host)
		username, _ = sect.ValueString("username", username)
		password, _ = sect.ValueString("password", password)
		port, _ = sect.ValueInt("port", port)
		vhost, _ = sect.ValueString("vhost", vhost)
		heartbeat, _ = sect.ValueInt("heartbeat", heartbeat)

		if exc, ok := sect.Section("exchange"); ok {
			exchangeName, _ = exc.ValueString("name", exchangeName)
			exchangeDurable, _ = exc.ValueBool("durable", exchangeDurable)
		}

		if qus, ok := sect.Section("queues"); ok {
			queuesAutodelete, _ = qus.ValueBool("autodelete", queuesAutodelete)
			queuesExclusive, _ = qus.ValueBool("exclusive", queuesExclusive)
		}

		jobQueueSize, _ = sect.ValueInt("jobqueue-size", jobQueueSize)
		workers, _ = sect.ValueInt("workers", workers)
	}

	// TODO: handle TLS
	acnl.cfg.connstr = util.AmqpUriBuilder(username, password, host, vhost, port, false)

	acnl.cfg.cfg = amqp.Config{
		Heartbeat: time.Duration(heartbeat) * time.Second,
	}

	acnl.cfg.exchange.name = exchangeName
	acnl.cfg.exchange.durable = exchangeDurable

	acnl.cfg.queues.autodelete = queuesAutodelete
	acnl.cfg.queues.exclusive = queuesExclusive

	acnl.sendQueue = make(chan *exchange, jobQueueSize)
	acnl.numOfWorkers = workers
	acnl.statusChan = make(chan channel.Status, workers)
	acnl.chanReqChan = make(chan struct{}, workers+1)
	acnl.newChanChan = make(chan struct {
		*amqp.Channel
		error
	})

	acnl.endpoint = &catalogue.Endpoint{
		Active:       true,
		Description:  acnl.cfg.vnfmDescr,
		Enabled:      true,
		Endpoint:     acnl.cfg.vnfmEndpoint,
		EndpointType: "RABBIT",
		Type:         acnl.cfg.vnfmType,
	}

	return acnl, acnl.spawn()
}
