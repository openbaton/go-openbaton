package catalogue

type LinkStatus string

const (
	LinkDown                           = LinkStatus("LINKDOWN")
	LinkNormalOperation                = LinkStatus("NORMALOPERATION")
	LinkDegradeOperation               = LinkStatus("DEGRADEDOPERATION")
	LinkOfflineThroughManagementAction = LinkStatus("OFFLINETHROUGHMANAGEMENTACTION")
)

// NetworkServiceRecord as defined by ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type NetworkServiceRecord struct {
	ID                       string                           `json:"id"`
	AutoScalePolicy          []*AutoScalePolicy               `json:"auto_scale_policy"`
	ConnectionPoint          []*ConnectionPoint               `json:"connection_point"`
	MonitoringParameter      []string                         `json:"monitoring_parameterid"`
	ServiceDeploymentFlavour NetworkServiceDeploymentFlavour  `json:"service_deployment_flavour"`
	Vendor                   string                           `json:"vendor"`
	ProjectID                string                           `json:"projectId"`
	Task                     string                           `json:"task"`
	Version                  string                           `json:"version"`
	VLR                      []*VirtualLinkRecord             `json:"vlr"`
	VNFR                     []*VirtualNetworkFunctionRecord  `json:"vnfr"`
	VNFDependency            []*VNFRecordDependency           `json:"vnf_dependency"`
	LifecycleEvent           []*LifecycleEvent                `json:"lifecycle_event"`
	VNFFGR                   []*VNFForwardingGraphRecord      `json:"vnffgr"`
	PNFR                     []*PhysicalNetworkFunctionRecord `json:"pnfr"`
	FaultManagementPolicy    []*FaultManagementPolicy         `json:"faultManagementPolicy"`
	DescriptorReference      string                           `json:"descriptor_reference"`
	ResourceReservation      string                           `json:"resource_reservation"`
	RuntimePolicyInfo        string                           `json:"runtime_policy_info"`
	Status                   Status                           `json:"status"`
	Notification             string                           `json:"notification"`
	LifecycleEventHistory    []*LifecycleEvent                `json:"lifecycle_event_history"`
	AuditLog                 string                           `json:"audit_log"`
	CreatedAt                string                           `json:"createdAt"`
	KeyNames                 []string                         `json:"keyNames"`
	Name                     string                           `json:"name"`
}

// PhysicalNetworkFunctionRecord based on ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type PhysicalNetworkFunctionRecord struct {
	ID                   string                      `json:"id"`
	Vendor               string                      `json:"vendor"`
	Version              string                      `json:"version"`
	Description          string                      `json:"description"`
	ProjectID            string                      `json:"projectId"`
	ConnectionPoint      []*ConnectionPoint          `json:"connection_point"`
	ParentNSID           string                      `json:"parent_ns_id"`
	DescriptorReference  string                      `json:"descriptor_reference"`
	VNFFGR               []*VNFForwardingGraphRecord `json:"vnffgr"`
	OamReference         string                      `json:"oam_reference"`
	ConnectedVirtualLink []*VirtualLinkRecord        `json:"connected_virtual_link"`
	PNFAddress           []string                    `json:"pnf_address"`
}

type Policy struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
}

type Status int

const (
	// Error
	StatusError Status = iota

	// Null
	StatusNull

	//Instantiated - Not Configured
	StatusInitialized

	// Inactive - Configured
	StatusInactive

	// Scaling
	StatusScaling

	// Active - Configured
	StatusActive

	// Terminated
	StatusTerminated

	// Resuming
	StatusResuming
)

type VirtualLinkRecord struct {
	Vendor                string                      `json:"vendor"`
	Version               string                      `json:"version"`
	NumberOfEndpoints     int                         `json:"number_of_endpoints"`
	ParentNs              string                      `json:"parent_ns"`
	VNFFGRReference       []*VNFForwardingGraphRecord `json:"vnffgr_reference"`
	DescriptorReference   string                      `json:"descriptor_reference"`
	VimID                 string                      `json:"vim_id"`
	AllocatedCapacity     []string                    `json:"allocated_capacity"`
	Status                LinkStatus                  `json:"status"`
	Notification          []string                    `json:"notification"`
	LifecycleEventHistory []*LifecycleEvent           `json:"lifecycle_event_history"`
	AuditLog              []string                    `json:"audit_log"`
	Connection            []string                    `json:"connection"`
}

// A VirtualNetworkFunctionRecord as described by ETSI GS NFV-MAN 001 V1.1.1
type VirtualNetworkFunctionRecord struct {
	ID                           string                   `json:"id"`
	HbVersion                    int                      `json:"hb_version"`
	AutoScalePolicy              []*AutoScalePolicy       `json:"auto_scale_policy"`
	ConnectionPoint              []*ConnectionPoint       `json:"connection_point"`
	ProjectID                    string                   `json:"projectId"`
	DeploymentFlavourKey         string                   `json:"deployment_flavour_key"`
	Configurations               *Configuration           `json:"configurations"`
	LifecycleEvent               []*LifecycleEvent        `json:"lifecycle_event"`
	LifecycleEventHistory        []*HistoryLifecycleEvent `json:"lifecycle_event_history"`
	Localization                 string                   `json:"localization"`
	MonitoringParameter          []string                 `json:"monitoring_parameter"`
	Vdu                          []*VirtualDeploymentUnit `json:"vdu"`
	Vendor                       string                   `json:"vendor"`
	Version                      string                   `json:"version"`
	VirtualLink                  []InternalVirtualLink    `json:"virtual_link"`
	ParentNsID                   string                   `json:"parent_ns_id"`
	DescriptorReference          string                   `json:"descriptor_reference"`
	VnfmID                       string                   `json:"vnfm_id"`
	ConnectedExternalVirtualLink []VirtualLinkRecord      `json:"connected_external_virtual_link"`
	VnfAddress                   []string                 `json:"vnf_address"`
	Status                       Status                   `json:"status"`
	Notification                 []string                 `json:"notification"`
	AuditLog                     string                   `json:"audit_log"`
	RuntimePolicyInfo            []string                 `json:"runtime_policy_info"`
	Name                         string                   `json:"name"`
	Type                         string                   `json:"type"`
	Endpoint                     string                   `json:"endpoint"`
	Task                         string                   `json:"task"`
	Requires                     *Configuration           `json:"requires"`
	Provides                     *Configuration           `json:"provides"`
	CyclicDependency             bool                     `json:"cyclic_dependency"`
	PackageID                    string                   `json:"packageId"`
}

type VNFCInstance struct {
	VimID        string        `json:"vim_id"`
	VcID         string        `json:"vc_id"`
	Hostname     string        `json:"hostname"`
	State        string        `json:"state"`
	VnfComponent *VNFComponent `json:"vnfComponent"`
	FloatingIps  []*Ip         `json:"floatingIps"`
	Ips          []*Ip         `json:"ips"`
}

// Based on ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type VNFForwardingGraphRecord struct {
	ID                    string                          `json:"id"`
	DescriptorReference   *VNFForwardingGraphDescriptor   `json:"descriptor_reference"`
	ParentNs              *NetworkServiceRecord           `json:"parent_ns"`
	DependentVirtualLink  []*VirtualLinkRecord            `json:"dependent_virtual_link"`
	Status                *Status                         `json:"status"`
	Notification          []string                        `json:"notification"`
	LifecycleEventHistory []*LifecycleEvent               `json:"lifecycle_event_history"`
	AuditLog              string                          `json:"audit_log"`
	NetworkForwardingPath *NetworkForwardingPath          `json:"network_forwarding_path"`
	ConnectionPoint       []*VNFDConnectionPoint          `json:"connection_point"`
	MemberVnfs            []*VirtualNetworkFunctionRecord `json:"member_vnfs"`
	Vendor                string                          `json:"vendor"`
	Version               string                          `json:"version"`
	NumberOfEndpoints     int                             `json:"number_of_endpoints"`
	NumberOfVnfs          int                             `json:"number_of_vnfs"`
	NumberOfPnfs          int                             `json:"number_of_pnfs"`
	NumberOfVirtualLinks  int                             `json:"number_of_virtual_links"`
}

type VNFRecordDependency struct {
	ID             string                               `json:"id"`
	Version        int                                  `json:"version"`
	Target         string                               `json:"target"`
	Parameters     map[string]*DependencyParameters     `json:"parameters"`
	VNFCParameters map[string]*VNFCDependencyParameters `json:"vnfcParameters"`
	IDType         map[string]string                    `json:"idType"`
}