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

// An extended Virtual Link based on ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type InternalVirtualLink struct {
	ID               string            `json:"id,omitempty"`
	HbVersion        int               `json:"hbVersion,omitempty"`
	ProjectID        string            `json:"projectId"`
	Shared           bool              `json:"shared,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
	ExtID            string            `json:"extId"`
	RootRequirement  string            `json:"root_requirement"`
	LeafRequirement  string            `json:"leaf_requirement"`
	QoS              []string          `json:"qos"`
	TestAccess       []string          `json:"test_access"`
	ConnectivityType []string          `json:"connectivity_type"`
	Name             string            `json:"name"`

	ConnectionPointsReferences []string `json:"connection_points_references"`
}

type NetworkForwardingPath struct {
	ID         string            `json:"id,omitempty"`
	HbVersion  int               `json:"hbVersion,omitempty"`
	ProjectID  string            `json:"projectId"`
	Shared     bool              `json:"shared,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Policy     *Policy           `json:"policy,omitempty"`
	Connection map[string]string `json:"connection,omitempty"`
}

type NFVEntityDescriptor struct {
	ID                        string                          `json:"id,omitempty"`
	HbVersion                 int                             `json:"hbVersion,omitempty"`
	ProjectID                 string                          `json:"projectId"`
	Shared                    bool                            `json:"shared,omitempty"`
	Metadata                  map[string]string               `json:"metadata,omitempty"`
	Name                      string                          `json:"name"`
	Vendor                    string                          `json:"vendor"`
	Version                   string                          `json:"version"`
	VNFFGDs                   []*VNFForwardingGraphDescriptor `json:"vnffgd"`
	VLDs                      []*VirtualLinkDescriptor        `json:"vld"`
	MonitoringParameters      []string                        `json:"monitoring_parameter"`
	ServiceDeploymentFlavours []*DeploymentFlavour            `json:"service_deployment_flavour"`
	AutoScalePolicies         []*AutoScalePolicy              `json:"auto_scale_policy"`
	ConnectionPoints          []*ConnectionPoint              `json:"connection_point"`
}

// VDUDepencency as described in ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type VDUDependency struct {
	ID        string                 `json:"id,omitempty"`
	HbVersion int                    `json:"hbVersion,omitempty"`
	ProjectID string                 `json:"projectId"`
	Shared    bool                   `json:"shared,omitempty"`
	Metadata  map[string]string      `json:"metadata,omitempty"`
	Source    *VirtualDeploymentUnit `json:"source,omitempty"`
	Target    *VirtualDeploymentUnit `json:"target,omitempty"`
}

type VirtualDeploymentUnit struct {
	ID                              string                     `json:"id,omitempty"`
	HbVersion                       int                        `json:"hbVersion,omitempty"`
	ProjectID                       string                     `json:"projectId"`
	Shared                          bool                       `json:"shared,omitempty"`
	Metadata                        map[string]string          `json:"metadata,omitempty"`
	Name                            string                     `json:"name"`
	VMImages                        []string                   `json:"vm_image"`
	ParentVDU                       string                     `json:"parent_vdu"`
	ComputationRequirement          string                     `json:"computation_requirement"`
	VirtualMemoryResourceElement    string                     `json:"virtual_memory_resource_element"`
	VirtualNetworkBandwidthResource string                     `json:"virtual_network_bandwidth_resource"`
	LifecycleEvents                 LifecycleEvents            `json:"lifecycle_event"`
	VduConstraint                   string                     `json:"vdu_constraint"`
	HighAvailability                *HighAvailability          `json:"high_availability,omitempty"`
	FaultManagementPolicies         []*VRFaultManagementPolicy `json:"fault_management_policy,omitempty"`
	ScaleInOut                      int                        `json:"scale_in_out"`
	VNFCs                           []*VNFComponent            `json:"vnfc"`
	VNFCInstances                   []*VNFCInstance            `json:"vnfc_instance"`
	MonitoringParameters            []string                   `json:"monitoring_parameter"`
	Hostname                        string                     `json:"hostname"`
	VIMInstanceNames                []string                   `json:"vimInstanceName"`
}

// VirtualLinkDescriptor as described in ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type VirtualLinkDescriptor struct {
	ID                string            `json:"id,omitempty"`
	HbVersion         int               `json:"hbVersion,omitempty"`
	ProjectID         string            `json:"projectId"`
	Shared            bool              `json:"shared,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	ExtID             string            `json:"extId"`
	RootRequirement   string            `json:"root_requirement"`
	LeafRequirement   string            `json:"leaf_requirement"`
	QoS               []string          `json:"qos"`
	TestAccess        []string          `json:"test_access"`
	ConnectivityType  []string          `json:"connectivity_type"`
	Name              string            `json:"name"`
	Vendor            string            `json:"vendor"`
	DescriptorVersion string            `json:"descriptor_version"`
	NumberOfEndpoints int               `json:"number_of_endpoints"`
	Connections       []string          `json:"connection"`
	VLDSecurity       *Security         `json:"vld_security,omitempty"`
}

// VirtualNetworkFunctionDescriptor as described in ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type VirtualNetworkFunctionDescriptor struct {
	ID                        string                          `json:"id,omitempty"`
	HbVersion                 int                             `json:"hbVersion,omitempty"`
	ProjectID                 string                          `json:"projectId"`
	Shared                    bool                            `json:"shared,omitempty"`
	Metadata                  map[string]string               `json:"metadata,omitempty"`
	Name                      string                          `json:"name"`
	Vendor                    string                          `json:"vendor"`
	Version                   string                          `json:"version"`
	VNFFGDs                   []*VNFForwardingGraphDescriptor `json:"vnffgd"`
	VLDs                      []*VirtualLinkDescriptor        `json:"vld"`
	MonitoringParameters      []string                        `json:"monitoring_parameter"`
	ServiceDeploymentFlavours []*DeploymentFlavour            `json:"service_deployment_flavour"`
	AutoScalePolicies         []*AutoScalePolicy              `json:"auto_scale_policy"`
	ConnectionPoints          []*ConnectionPoint              `json:"connection_point"`
	LifecycleEvents           LifecycleEvents                 `json:"lifecycle_event"`
	Configurations            *Configuration                  `json:"configurations,omitempty"`
	VDUs                      []*VirtualDeploymentUnit        `json:"vdu"`
	VirtualLinks              []*InternalVirtualLink          `json:"virtual_link"`
	VDUDependencies           []*VDUDependency                `json:"vdu_dependency"`
	DeploymentFlavours        []*VNFDeploymentFlavour         `json:"deployment_flavour"`
	ManifestFile              string                          `json:"manifest_file"`
	ManifestFileSecurity      []*Security                     `json:"manifest_file_security"`
	Type                      string                          `json:"type"`
	Endpoint                  string                          `json:"endpoint"`
	VNFPackageLocation        string                          `json:"vnfPackageLocation"`
	Requires                  map[string]*RequiresParameters  `json:"requires,omitempty"`
	Provides                  []string                        `json:"provides,omitempty"`
	CyclicDependency          bool                            `json:"cyclicDependency"`
	VNFDConnectionPoints      []*VNFDConnectionPoint          `json:"VNFDConnection_point"`
}

// A VNFComponent as defined by ETSI GS NFV-MAN 001 V1.1.1
type VNFComponent struct {
	ID               string                 `json:"id,omitempty"`
	HbVersion        int                    `json:"hbVersion,omitempty"`
	ProjectID        string                 `json:"projectId"`
	Shared           bool                   `json:"shared,omitempty"`
	Metadata         map[string]string      `json:"metadata,omitempty"`
	ConnectionPoints []*VNFDConnectionPoint `json:"connection_point"`
}

// Virtual Network Function Descriptor Connection Point as defined by
// ETSI GS NFV-MAN 001 V1.1.1
type VNFDConnectionPoint struct {
	ID         string            `json:"id,omitempty"`
	HbVersion  int               `json:"hbVersion,omitempty"`
	ProjectID  string            `json:"projectId"`
	Shared     bool              `json:"shared,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Type       string            `json:"type"`
	FixedIp    string            `json:"floatingIp"`
	ChosenPool string            `json:"chosenPool"`

	VirtualLinkReference string `json:"virtual_link_reference"`
	FloatingIP           string `json:"floatingIp"`
	InterfaceID          int    `json:"interfaceId"`
}

// VNFForwardingGraphDescriptor as defined by ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type VNFForwardingGraphDescriptor struct {
	ID                     string                   `json:"id,omitempty"`
	HbVersion              int                      `json:"hbVersion,omitempty"`
	ProjectID              string                   `json:"projectId"`
	Shared                 bool                     `json:"shared,omitempty"`
	Metadata               map[string]string        `json:"metadata,omitempty"`
	Vendor                 string                   `json:"vendor"`
	Version                string                   `json:"version"`
	NumberOfEndpoints      int                      `json:"number_of_endpoints"`
	NumberOfVirtualLinks   int                      `json:"number_of_virtual_links"`
	DependentVirtualLinks  []*VirtualLinkDescriptor `json:"dependent_virtual_link"`
	NetworkForwardingPaths []*NetworkForwardingPath `json:"network_forwarding_path"`
	ConnectionPoints       []*ConnectionPoint       `json:"connection_point"`
	DescriptorVersion      string                   `json:"descriptor_version"`
	ConstituentVnfs        []*ConstituentVNF        `json:"constituent_vnfs"`
	VnffgdSecurity         *Security                `json:"vnffgd_security,omitempty"`
}
