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
	"github.com/openbaton/go-openbaton/vnfm"
	"github.com/openbaton/go-openbaton/vnfm/channel"
	"github.com/openbaton/go-openbaton/vnfm/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	vnfm.Register("amqp", amqpDriver{})
}

type amqpDriver struct{}

func (amqpDriver) Init(cnf *config.Config, log *log.Logger) (channel.Channel, error) {
	ret, err := newChannel(cnf, log)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
