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
)

type plugError struct {
	Message string `json:"detailMessage"`
}

func (e plugError) Error() string {
	return e.Message
}

type request struct {
	MethodName string            `json:"methodName"`
	Parameters []json.RawMessage `json:"parameters"`
}

type response struct {
	Answer    interface{} `json:"answer,omitempty"`
	Exception error       `json:"exception,omitempty"`
}
