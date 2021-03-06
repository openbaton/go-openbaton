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

package messages

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fatih/structs"
	"github.com/openbaton/go-openbaton/catalogue"
)

var (
	ErrMalformedMessage = errors.New("received JSON is malformed")
)

func Marshal(msg NFVMessage) ([]byte, error) {
	return json.MarshalIndent(msg, "", "  ")
}

func Unmarshal(msgBytes []byte, from SenderType) (NFVMessage, error) {
	msg := message{from: from}
	err := json.Unmarshal(msgBytes, &msg)

	return &msg, err
}

func UnmarshalCredentials(msgBytes []byte) (*catalogue.ManagerCredentials, error) {
	res := catalogue.ManagerCredentials{}

	err := json.Unmarshal(msgBytes, &res)

	return &res, err
}

func (msg *message) MarshalJSON() ([]byte, error) {
	// We need to serialize the message structure into a
	// compatible JSON message.
	// To achieve this, it's necessary to generate a temporary map
	// in which inject the necessary fields into before any serialization
	// of the message can occour.
	s := structs.New(msg.Content())

	// This line below tells structs to use the value contained in the "json" tags as keys
	// of the generated map.
	s.TagName = "json"

	tmpMap := s.Map()

	// Inject the "action" field into the map

	tmpMap["action"] = msg.Action()

	return json.Marshal(tmpMap)
}

func (msg *message) UnmarshalJSON(data []byte) error {
	// From should be already set before calling this function!

	// The action field must be extracted from the received object.
	// Using json.RawMessages reduces considerably the cost of unmarshalling the data
	// twice.
	tmp := make(map[string]json.RawMessage)

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if rawAction, ok := tmp["action"]; ok {
		var action catalogue.Action
		if err := json.Unmarshal(rawAction, &action); err != nil {
			return err
		}

		msg.action = action
	} else {
		return ErrMalformedMessage
	}

	switch msg.From() {
	case NFVO:
		return msg.unmarshalNFVOMessage(data)

	case VNFM:
		return msg.unmarshalVNFMMessage(data)

	default:
		return fmt.Errorf("invalid sender type %v", msg.From())
	}
}

func (msg *message) unmarshalNFVOMessage(data []byte) error {
	switch msg.Action() {
	case catalogue.ActionGrantOperation:
		msg.content = &OrGrantLifecycleOperation{}

	case catalogue.ActionScaleIn:
		fallthrough
	case catalogue.ActionScaleOut:
		fallthrough
	case catalogue.ActionScaling:
		msg.content = &OrScaling{}

	case catalogue.ActionError:
		msg.content = &OrError{}

	case catalogue.ActionInstantiate:
		tmp := make(map[string]json.RawMessage)

		if err := json.Unmarshal(data, &tmp); err != nil {
			return err
		}
		tmpContent := &OrInstantiate{}
		json.Unmarshal(data, tmpContent)
		tmpContent.VIMInstances = make(map[string][]interface{})
		tmpVim := make(map[string][]json.RawMessage)
		if err := json.Unmarshal(tmp["vimInstances"], &tmpVim); err != nil {
			return err
		}
		for k, v := range tmpVim { //map[string][]interface{}
			tmpContent.VIMInstances[k] = make([]interface{}, len(v))
			for i, vimRaw := range v {
				tmpContent.VIMInstances[k][i] = parseVim(vimRaw)
			}
		}

		msg.content = tmpContent
		return nil

	case catalogue.ActionHeal:
		msg.content = &OrHealVNFRequest{}

	case catalogue.ActionUpdate:
		msg.content = &OrUpdate{}

	case catalogue.ActionStart:
		fallthrough
	case catalogue.ActionStop:
		msg.content = &OrStartStop{}

	default:
		msg.content = &OrGeneric{}
	}

	return json.Unmarshal(data, &msg.content)
}
func parseVim(rawVim json.RawMessage) interface{} {
	tmp := make(map[string]json.RawMessage)

	if err := json.Unmarshal(rawVim, &tmp); err != nil {
		return err
	}

	var _typ string
	if err := json.Unmarshal(tmp["type"], &_typ); err != nil {
		return err
	}
	if _typ == "docker" {
		res := &catalogue.DockerVimInstance{}
		if err := json.Unmarshal(rawVim, res); err != nil {
			return err
		}
		return res
	} else if string(tmp["type"]) == "openstack" {
		res := &catalogue.OpenstackVimInstance{}
		if err := json.Unmarshal(rawVim, res); err != nil {
			return err
		}
		return res
	} else {
		res := &catalogue.BaseVimInstance{}
		if err := json.Unmarshal(rawVim, res); err != nil {
			return err
		}
		return res
	}
}

func (msg *message) unmarshalVNFMMessage(data []byte) error {
	switch msg.Action() {
	case catalogue.ActionAllocateResources:
		msg.content = &VNFMAllocateResources{}

	case catalogue.ActionError:
		msg.content = &VNFMError{}

	case catalogue.ActionInstantiate:
		msg.content = &VNFMInstantiate{}

	case catalogue.ActionScaled:
		msg.content = &VNFMScaled{}

	case catalogue.ActionScaling:
		msg.content = &VNFMScaling{}

	case catalogue.ActionHeal:
		msg.content = &VNFMHealed{}

	case catalogue.ActionStart:
		msg.content = &VNFMStartStop{}

	case catalogue.ActionStop:
		msg.content = &VNFMStartStop{}

	default:
		msg.content = &VNFMGeneric{}
	}

	return json.Unmarshal(data, &msg.content)
}

func sanitizeAction(action catalogue.Action) catalogue.Action {
	switch action {
	// Ignore valid actions.
	// Convert invalid ones into catalogue.NoActionSpecified
	case catalogue.ActionGrantOperation:
	case catalogue.ActionAllocateResources:
	case catalogue.ActionScaleIn:
	case catalogue.ActionScaleOut:
	case catalogue.ActionScaling:
	case catalogue.ActionError:
	case catalogue.ActionReleaseResources:
	case catalogue.ActionInstantiate:
	case catalogue.ActionModify:
	case catalogue.ActionHeal:
	case catalogue.ActionUpdateVNFR:
	case catalogue.ActionUpdate:
	case catalogue.ActionScaled:
	case catalogue.ActionReleaseResourcesFinish:
	case catalogue.ActionInstantiateFinish:
	case catalogue.ActionConfigure:
	case catalogue.ActionStart:
	case catalogue.ActionStop:
	case catalogue.ActionResume:

	default:
		return catalogue.NoActionSpecified
	}

	return action
}
