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
	"github.com/openbaton/go-openbaton/catalogue"
)

type orMessage struct{}

func (orMessage) From() SenderType {
	return NFVO
}

type OrAllocateResources struct {
	orMessage

	VDUSet []*catalogue.VirtualDeploymentUnit `json:"vduSet,omitempty"`
}

func (OrAllocateResources) DefaultAction() catalogue.Action {
	return catalogue.ActionAllocateResources
}

type OrError struct {
	orMessage

	VNFR    *catalogue.VirtualNetworkFunctionRecord `json:"vnfr,omitempty"`
	Message string                                  `json:"message,omitempty"`
}

func (OrError) DefaultAction() catalogue.Action {
	return catalogue.ActionError
}

type OrGeneric struct {
	orMessage

	VNFR           *catalogue.VirtualNetworkFunctionRecord `json:"vnfr,omitempty"`
	VNFRDependency *catalogue.VNFRecordDependency          `json:"vnfrd,omitempty"`
}

func (OrGeneric) DefaultAction() catalogue.Action {
	return catalogue.NoActionSpecified
}

type OrGrantLifecycleOperation struct {
	orMessage

	GrantAllowed bool                                    `json:"grantAllowed"`
	VDUVIM       map[string]*catalogue.BaseVimInstance   `json:"vduVim,omitempty"`
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
}

func (OrGrantLifecycleOperation) DefaultAction() catalogue.Action {
	return catalogue.ActionGrantOperation
}

type OrHealVNFRequest struct {
	orMessage

	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance,omitempty"`
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	Cause        string                                  `json:"cause,omitempty"`
}

func (OrHealVNFRequest) DefaultAction() catalogue.Action {
	return catalogue.ActionHeal
}

type OrInstantiate struct {
	orMessage

	VNFD            *catalogue.VirtualNetworkFunctionDescriptor `json:"vnfd,omitempty"`
	VNFDFlavour     *catalogue.VNFDeploymentFlavour             `json:"vnfdf,omitempty"`
	VNFInstanceName string                                      `json:"vnfInstanceName,omitempty"`
	VLRs            []*catalogue.VirtualLinkRecord              `json:"vlrs,omitempty"`
	Extension       map[string]string                           `json:"extension,omitempty"`
	VIMInstances    map[string][]*catalogue.BaseVimInstance     `json:"vimInstances,omitempty"`
	VNFPackage      *catalogue.VNFPackage                       `json:"vnfPackage,omitempty"`
	Keys            []*catalogue.Key                            `json:"keys,omitempty"`
}

func (OrInstantiate) DefaultAction() catalogue.Action {
	return catalogue.ActionInstantiate
}

type OrScaling struct {
	orMessage

	Component    *catalogue.VNFComponent                 `json:"component,omitempty"`
	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance,omitempty"`
	VIMInstance  *catalogue.BaseVimInstance              `json:"vimInstance,omitempty"`
	VNFPackage   *catalogue.VNFPackage                   `json:"vnfPackage,omitempty"`
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	Dependency   *catalogue.VNFRecordDependency          `json:"dependency,omitempty"`
	Mode         string                                  `json:"mode,omitempty"`
	Extension    map[string]string                       `json:"extension,omitempty"`
}

func (OrScaling) DefaultAction() catalogue.Action {
	return catalogue.ActionScaling
}

type OrStartStop struct {
	orMessage

	VNFR           *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	VNFCInstance   *catalogue.VNFCInstance                 `json:"vnfcInstance,omitempty"`
	VNFRDependency *catalogue.VNFRecordDependency          `json:"vnfrDependency,omitempty"`
}

func (OrStartStop) DefaultAction() catalogue.Action {
	return catalogue.ActionStart
}

type OrUpdate struct {
	orMessage

	Script *catalogue.Script                       `json:"script,omitempty"`
	VNFR   *catalogue.VirtualNetworkFunctionRecord `json:"vnfr,omitempty"`
}

func (OrUpdate) DefaultAction() catalogue.Action {
	return catalogue.ActionUpdate
}
