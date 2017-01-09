package catalogue

type LinkStatus string

const (
	LinkDown                           = LinkStatus("LINKDOWN")
	LinkNormalOperation                = LinkStatus("NORMALOPERATION")
	LinkDegradeOperation               = LinkStatus("DEGRADEDOPERATION")
	LinkOfflineThroughManagementAction = LinkStatus("OFFLINETHROUGHMANAGEMENTACTION")
)

// Component represents a generic component.
type Component interface{}

// NetworkServiceRecord as defined by ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type NetworkServiceRecord struct {
	ID                       string                           `json:"id,omitempty"`
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
	ID                   string                      `json:"id,omitempty"`
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
	ID      string `json:"id,omitempty"`
	Version int    `json:"version"`
}

type Status string

const (
	// Error
	StatusError = Status("ERROR")

	// Null
	StatusNull = Status("NULL")

	//Instantiated - Not Configured
	StatusInitialized = Status("INITIALIZED")

	// Inactive - Configured
	StatusInactive = Status("INACTIVE")

	// Scaling
	StatusScaling = Status("SCALING")

	// Active - Configured
	StatusActive = Status("ACTIVE")

	// Terminated
	StatusTerminated = Status("TERMINATED")

	// Resuming
	StatusResuming = Status("RESUMING")
)

type VirtualLinkRecord struct {
	VirtualLink

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

type VNFCInstance struct {
	VNFComponent

	VIMID              string        `json:"vim_id"`
	VCID               string        `json:"vc_id"`
	Hostname           string        `json:"hostname"`
	State              string        `json:"state"`
	NestedVNFComponent *VNFComponent `json:"vnfComponent,omitempty"`
	FloatingIPs        []*Ip         `json:"floatingIps"`
	IPs                []*Ip         `json:"ips"`
}

// Based on ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type VNFForwardingGraphRecord struct {
	ID                    string                          `json:"id,omitempty"`
	DescriptorReference   *VNFForwardingGraphDescriptor   `json:"descriptor_reference"`
	ParentNS              *NetworkServiceRecord           `json:"parent_ns"`
	DependentVirtualLink  []*VirtualLinkRecord            `json:"dependent_virtual_link"`
	Status                *Status                         `json:"status,omitempty"`
	Notification          []string                        `json:"notification"`
	LifecycleEventHistory []*LifecycleEvent               `json:"lifecycle_event_history"`
	AuditLog              string                          `json:"audit_log"`
	NetworkForwardingPath *NetworkForwardingPath          `json:"network_forwarding_path,omitempty"`
	ConnectionPoint       []*VNFDConnectionPoint          `json:"connection_point"`
	MemberVNFs            []*VirtualNetworkFunctionRecord `json:"member_vnfs"`
	Vendor                string                          `json:"vendor"`
	Version               string                          `json:"version"`
	NumberOfEndpoints     int                             `json:"number_of_endpoints"`
	NumberOfVNFs          int                             `json:"number_of_vnfs"`
	NumberOfPNFs          int                             `json:"number_of_pnfs"`
	NumberOfVirtualLinks  int                             `json:"number_of_virtual_links"`
}

type VNFRecordDependency struct {
	ID             string                               `json:"id,omitempty"`
	Version        int                                  `json:"version"`
	Target         string                               `json:"target"`
	Parameters     map[string]*DependencyParameters     `json:"parameters"`
	VNFCParameters map[string]*VNFCDependencyParameters `json:"vnfcParameters"`
	IDType         map[string]string                    `json:"idType"`
}
