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

package catalogue

type Criteria struct {
	ID                 string       `json:"id,omitempty"`
	Version            int          `json:"version"`
	Name               string       `json:"name"`
	ParameterRef       string       `json:"parameter_ref"`
	Function           string       `json:"function"`
	VNFCSelector       VNFCSelector `json:"vnfc_selector"`
	ComparisonOperator string       `json:"comparison_operator"`
	Threshold          string       `json:"threshold"`
}

type FaultManagementAction string

const (
	FaultRestart              = FaultManagementAction("RESTART")
	FaultReinstantiateService = FaultManagementAction("REINSTANTIATE_SERVICE")
	FaultHeal                 = FaultManagementAction("HEAL")
	FaultReinstantiate        = FaultManagementAction("REINSTANTIATE")
	FaultSwitchToStandby      = FaultManagementAction("SWITCH_TO_STANDBY")
	FaultSwitchToActive       = FaultManagementAction("SWITCH_TO_ACTIVE")
)

type FaultManagementPolicy struct {
	ID       string            `json:"id,omitempty"`
	Name     string            `json:"name"`
	VNFAlarm bool              `json:"VNFAlarm"`
	Period   int               `json:"period"`
	Severity PerceivedSeverity `json:"severity"`
	Criteria []*Criteria       `json:"criteria"`
	Version  int               `json:"version"`
}

type VNFCSelector string

const (
	SelectorAtLeastOne = VNFCSelector("at_least_one")
	SelectorAll        = VNFCSelector("all")
)

type VRFaultManagementPolicy struct {
	FaultManagementPolicy

	Action FaultManagementAction `json:"action"`
}
