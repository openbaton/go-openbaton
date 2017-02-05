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

package util

import (
	"bytes"
	"runtime"
	"strconv"
)

const (
	// UnknownFunc is a string representing an unknown function.
	UnknownFunc = "<unknown function>"
)

func AmqpUriBuilder(username, password, host, vhost string, port int, tls bool) string {
	buffer := bytes.NewBufferString("amqp")

	if tls {
		buffer.WriteRune('s')
	}

	buffer.WriteString("://")
	if username != "" {
		buffer.WriteString(username)

		if password != "" {
			buffer.WriteRune(':')
			buffer.WriteString(password)
		}

		buffer.WriteRune('@')
	}

	if host != "" {
		buffer.WriteString(host)
	}

	if port > 0 {
		buffer.WriteRune(':')
		buffer.WriteString(strconv.Itoa(port))
	}

	if vhost != "" {
		buffer.WriteRune('/')
		buffer.WriteString(vhost)
	}

	return buffer.String()
}

func FuncName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return UnknownFunc
	}

	if f := runtime.FuncForPC(pc); f != nil {
		return f.Name()
	}

	return UnknownFunc
}
