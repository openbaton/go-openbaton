package openbaton

// An extended Virtual Link based on ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type InternalVirtualLink struct {
	VirtualLink
	ConnectionPointsReferences []string `json:"connection_points_references"`
}

type NetworkForwardingPath struct {
	ID         string            `json:"id"`
	Version    int               `json:"version"`
	Policy     *Policy           `json:"policy"`
	Connection map[string]string `json:"connection"`
}

type VirtualDeploymentUnit struct {
	ID                              string                     `json:"id"`
	Version                         int                        `json:"version"`
	ProjectID                       string                     `json:"projectId"`
	Name                            string                     `json:"name"`
	VmImage                         []string                   `json:"vm_image"`
	ParentVdu                       string                     `json:"parent_vdu"`
	ComputationRequirement          string                     `json:"computation_requirement"`
	VirtualMemoryResourceElement    string                     `json:"virtual_memory_resource_element"`
	VirtualNetworkBandwidthResource string                     `json:"virtual_network_bandwidth_resource"`
	LifecycleEvent                  []*LifecycleEvent          `json:"lifecycle_event"`
	VduConstraint                   string                     `json:"vdu_constraint"`
	HighAvailability                *HighAvailability          `json:"high_availability"`
	FaultManagementPolicy           []*VRFaultManagementPolicy `json:"fault_management_policy"`
	ScaleInOut                      int                        `json:"scale_in_out"`
	Vnfc                            []*VNFComponent            `json:"vnfc"`
	VnfcInstance                    []*VNFCInstance            `json:"vnfc_instance"`
	MonitoringParameter             []string                   `json:"monitoring_parameter"`
	Hostname                        string                     `json:"hostname"`
	VimInstanceName                 []string                   `json:"vimInstanceName"`
}

// VirtualLinkDescriptor as described in ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type VirtualLinkDescriptor struct {
	VirtualLink
	ProjectID         string    `json:"projectId"`
	Vendor            string    `json:"vendor"`
	DescriptorVersion string    `json:"descriptor_version"`
	NumberOfEndpoints int       `json:"number_of_endpoints"`
	Connection        []string  `json:"connection"`
	VLDSecurity       *Security `json:"vld_security"`
	Name              string    `json:"name"`
}

// A Virtual Network Function Component as defined by ETSI GS NFV-MAN 001 V1.1.1
type VNFComponent struct {
	ID              string                 `json:"id"`
	Version         int                    `json:"version"`
	ConnectionPoint []*VNFDConnectionPoint `json:"connection_component"`
}

// Virtual Network Function Descriptor Connection Point as defined by
// ETSI GS NFV-MAN 001 V1.1.1
type VNFDConnectionPoint struct {
	VirtualLinkReference string `json:"virtual_link_reference"`
	FloatingIp           string `json:"floatingIp"`
}

// VNFForwardingGraphDescriptor as defined by ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type VNFForwardingGraphDescriptor struct {
	ID                    string                   `json:"id"`
	HbVersion             int                      `json:"hb_version"`
	Vendor                string                   `json:"vendor"`
	Version               string                   `json:"version"`
	NumberOfEndpoints     int                      `json:"number_of_endpoints"`
	NumberOfVirtualLinks  int                      `json:"number_of_virtual_links"`
	DependentVirtualLink  []*VirtualLinkDescriptor `json:"dependent_virtual_link"`
	NetworkForwardingPath []*NetworkForwardingPath `json:"network_forwarding_path"`
	ConnectionPoint       []*ConnectionPoint       `json:"connection_point"`
	DescriptorVersion     string                   `json:"descriptor_version"`
	ConstituentVnfs       []*ConstituentVNF        `json:"constituent_vnfs"`
	VnffgdSecurity        *Security                `json:"vnffgd_security"`
}
