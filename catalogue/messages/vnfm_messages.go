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

import "github.com/openbaton/go-openbaton/catalogue"

type baseMessage struct{}

type vnfmMessage baseMessage

func (vnfmMessage) From() SenderType {
	return VNFM
}

type VNFMAllocateResources struct {
	vnfmMessage

	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	VIMInstances map[string]*catalogue.VIMInstance       `json:"vimInstances,omitempty"`
	Userdata     string                                  `json:"userdata,omitempty"`
	KeyPairs     []*catalogue.Key                        `json:"keyPairs,omitempty"`
}

func (VNFMAllocateResources) DefaultAction() catalogue.Action {
	return catalogue.ActionAllocateResources
}

type VNFMError struct {
	vnfmMessage

	NSRID     string                                  `json:"nsrId,omitempty"`
	VNFR      *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	Exception map[string]interface{}                  `json:"exception,omitempty"` // I don't know how to deserialize a Java exception
}

func (VNFMError) DefaultAction() catalogue.Action {
	return catalogue.ActionError
}

type VNFMGeneric struct {
	vnfmMessage

	VNFR                *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	VNFRecordDependency *catalogue.VNFRecordDependency          `json:"vnfRecordDependency,omitempty"`
}

func (VNFMGeneric) DefaultAction() catalogue.Action {
	return catalogue.NoActionSpecified
}

type VNFMGrantLifecycleOperation struct {
	vnfmMessage

	VNFD                 *catalogue.VirtualNetworkFunctionDescriptor `json:"virtualNetworkFunctionDescriptor,omitempty"`
	VDUSet               []*catalogue.VirtualDeploymentUnit          `json:"vduSet,omitempty"`
	DeploymentFlavourKey string                                      `json:"deploymentFlavourKey,omitempty"`
}

func (VNFMGrantLifecycleOperation) DefaultAction() catalogue.Action {
	return catalogue.ActionGrantOperation
}

type VNFMHealed struct {
	vnfmMessage

	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance,omitempty"`
	Cause        string                                  `json:"cause,omitempty"`
}

func (VNFMHealed) DefaultAction() catalogue.Action {
	return catalogue.ActionHeal
}

type VNFMInstantiate struct {
	vnfmMessage

	VNFR *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
}

func (VNFMInstantiate) DefaultAction() catalogue.Action {
	return catalogue.ActionInstantiate
}

type VNFMScaled struct {
	vnfmMessage

	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance,omitempty"`
}

func (VNFMScaled) DefaultAction() catalogue.Action {
	return catalogue.ActionScaled
}

type VNFMScaling struct {
	vnfmMessage

	VNFR     *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	UserData string                                  `json:"userData,omitempty"`
}

func (VNFMScaling) DefaultAction() catalogue.Action {
	return catalogue.ActionScaling
}

type VNFMStartStop struct {
	vnfmMessage

	VNFR           *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	VNFCInstance   *catalogue.VNFCInstance                 `json:"vnfcInstance,omitempty"`
	VNFRDependency *catalogue.VNFRecordDependency          `json:"vnfrDependency,omitempty"`
}

func (VNFMStartStop) DefaultAction() catalogue.Action {
	return catalogue.ActionStart
}




