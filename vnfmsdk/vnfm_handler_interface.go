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

package vnfmsdk

import (
	"github.com/openbaton/go-openbaton/catalogue"
)



// The Handler interface defines an abstraction of the operations that a VNFM should provide.
type HandlerVnfm interface {
	// ActionForResume uses the given VNFR and VNFCInstance to return a valid
	// action for resume. NoSuchAction is returned in case no such Action exists.
	ActionForResume(vnfr *catalogue.VirtualNetworkFunctionRecord,
		vnfcInstance *catalogue.VNFCInstance) catalogue.Action

	// CheckInstantiationFeasibility allows the VNFM to verify if the VNF instantiation is possible.
	CheckInstantiationFeasibility() error

	Configure(vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error)

	HandleError(vnfr *catalogue.VirtualNetworkFunctionRecord) error

	Heal(vnfr *catalogue.VirtualNetworkFunctionRecord,
		component *catalogue.VNFCInstance, cause string) (*catalogue.VirtualNetworkFunctionRecord, error)

	// Instantiate allows to create a VNF instance.
	Instantiate(vnfr *catalogue.VirtualNetworkFunctionRecord, scripts interface{},
		vimInstances map[string][]*catalogue.VIMInstance) (*catalogue.VirtualNetworkFunctionRecord, error)

	// Modify allows making structural changes (e.g.configuration, topology, behavior, redundancy model) to a VNF instance.
	Modify(vnfr *catalogue.VirtualNetworkFunctionRecord,
		dependency *catalogue.VNFRecordDependency) (*catalogue.VirtualNetworkFunctionRecord, error)

	// Query allows retrieving a VNF instance state and attributes. (not implemented)
	Query() error

	Resume(vnfr *catalogue.VirtualNetworkFunctionRecord,
		vnfcInstance *catalogue.VNFCInstance,
		dependency *catalogue.VNFRecordDependency) (*catalogue.VirtualNetworkFunctionRecord, error)

	// Scale allows scaling (out / in, up / down) a VNF instance.
	Scale(scaleInOrOut catalogue.Action,
		vnfr *catalogue.VirtualNetworkFunctionRecord,
		component catalogue.Component,
		scripts interface{},
		dependency *catalogue.VNFRecordDependency) (*catalogue.VirtualNetworkFunctionRecord, error)

	// Start starts a VNFR.
	Start(vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error)

	StartVNFCInstance(vnfr *catalogue.VirtualNetworkFunctionRecord,
		vnfcInstance *catalogue.VNFCInstance) (*catalogue.VirtualNetworkFunctionRecord, error)

	// Stop stops a previously created VNF instance.
	Stop(vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error)

	StopVNFCInstance(vnfr *catalogue.VirtualNetworkFunctionRecord,
		vnfcInstance *catalogue.VNFCInstance) (*catalogue.VirtualNetworkFunctionRecord, error)

	// Terminate allows terminating gracefully or forcefully a previously created VNF instance.
	Terminate(vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error)

	// UpdateSoftware allows applying a minor / limited software update(e.g.patch) to a VNF instance.
	UpdateSoftware(script *catalogue.Script,
		vnfr *catalogue.VirtualNetworkFunctionRecord) (*catalogue.VirtualNetworkFunctionRecord, error)

	// UpgradeSoftware allows deploying a new software release to a VNF instance.
	UpgradeSoftware() error

	// UserData returns a string containing UserData.
	UserData() string

}
