package catalogue

// An extended Virtual Link based on ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type InternalVirtualLink struct {
	VirtualLink

	ConnectionPointsReferences []string `json:"connection_points_references"`
}

type NetworkForwardingPath struct {
	ID         string            `json:"id,omitempty"`
	Version    int               `json:"version"`
	Policy     *Policy           `json:"policy,omitempty"`
	Connection map[string]string `json:"connection,omitempty"`
}

type NFVEntityDescriptor struct {
	ID                        string                          `json:"id,omitempty"`
	HbVersion                 int                             `json:"hb_version,omitempty"`
	Name                      string                          `json:"name"`
	ProjectID                 string                          `json:"projectId"`
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
	ID      string                 `json:"id,omitempty"`
	Version int                    `json:"version"`
	Source  *VirtualDeploymentUnit `json:"source,omitempty"`
	Target  *VirtualDeploymentUnit `json:"target,omitempty"`
}

type VirtualDeploymentUnit struct {
	ID                              string                     `json:"id,omitempty"`
	Version                         int                        `json:"version"`
	ProjectID                       string                     `json:"projectId"`
	Name                            string                     `json:"name"`
	VMImages                        []string                   `json:"vm_image"`
	ParentVDU                       string                     `json:"parent_vdu"`
	ComputationRequirement          string                     `json:"computation_requirement"`
	VirtualMemoryResourceElement    string                     `json:"virtual_memory_resource_element"`
	VirtualNetworkBandwidthResource string                     `json:"virtual_network_bandwidth_resource"`
	LifecycleEvents                 []*LifecycleEvent          `json:"lifecycle_event"`
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
	VirtualLink
	ProjectID         string    `json:"projectId"`
	Vendor            string    `json:"vendor"`
	DescriptorVersion string    `json:"descriptor_version"`
	NumberOfEndpoints int       `json:"number_of_endpoints"`
	Connections       []string  `json:"connection"`
	VLDSecurity       *Security `json:"vld_security,omitempty"`
	Name              string    `json:"name"`
}

// VirtualNetworkFunctionDescriptor as described in ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type VirtualNetworkFunctionDescriptor struct {
	NFVEntityDescriptor

	LifecycleEvents      []*LifecycleEvent              `json:"lifecycle_event"`
	Configurations       *Configuration                 `json:"configurations,omitempty"`
	VDUs                 []*VirtualDeploymentUnit       `json:"vdu"`
	VirtualLinks         []*InternalVirtualLink         `json:"virtual_link"`
	VDUDependencies      []*VDUDependency               `json:"vdu_dependency"`
	DeploymentFlavours   []*VNFDeploymentFlavour        `json:"deployment_flavour"`
	ManifestFile         string                         `json:"manifest_file"`
	ManifestFileSecurity []*Security                    `json:"manifest_file_security"`
	Type                 string                         `json:"type"`
	Endpoint             string                         `json:"endpoint"`
	VNFPackageLocation   string                         `json:"vnfPackageLocation"`
	Requires             map[string]*RequiresParameters `json:"requires,omitempty"`
	Provides             []string                       `json:"provides,omitempty"`
	CyclicDependency     bool                           `json:"cyclicDependency"`
	ConnectionPoints     []*ConnectionPoint             `json:"connection_point"`
	VNFDConnectionPoints []*VNFDConnectionPoint         `json:"VNFDConnection_point"`
}

// A VNFComponent as defined by ETSI GS NFV-MAN 001 V1.1.1
type VNFComponent struct {
	ID               string                 `json:"id,omitempty"`
	Version          int                    `json:"version"`
	ConnectionPoints []*VNFDConnectionPoint `json:"connection_component"`
}

// Virtual Network Function Descriptor Connection Point as defined by
// ETSI GS NFV-MAN 001 V1.1.1
type VNFDConnectionPoint struct {
	ConnectionPoint

	VirtualLinkReference string `json:"virtual_link_reference"`
	FloatingIP           string `json:"floatingIp"`
}

// VNFForwardingGraphDescriptor as defined by ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type VNFForwardingGraphDescriptor struct {
	ID                     string                   `json:"id,omitempty"`
	HbVersion              int                      `json:"hb_version"`
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

