package openbaton

// NetworkServiceRecord as defined by ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type NetworkServiceRecord struct {
	ID string `json:"id"`
	AutoScalePolicy []*AutoScalePolicy `json:"auto_scale_policy"`
	ConnectionPoint []*ConnectionPoint `json:"connection_point"`
	MonitoringParameter []string `json:"monitoring_parameterid"`
	ServiceDeploymentFlavour NetworkServiceDeploymentFlavour `json:"service_deployment_flavour"`
	Vendor string `json:"vendor"`
	ProjectID string `json:"projectId"`
	Task string `json:"task"`
	Version string `json:"version"`
	VLR []*VirtualLinkRecord `json:"vlr"`
	VNFR []*VirtualNetworkFunctionRecord `json:"vnfr"`
	VNFDependency []*VNFRecordDependency `json:"vnf_dependency"`
	LifecycleEvent []*LifecycleEvent `json:"lifecycle_event"`
	VNFFGR []*VNFForwardingGraphRecord `json:"vnffgr"`
	PNFR []*PhysicalNetworkFunctionRecord `json:"pnfr"`
	FaultManagementPolicy []*FaultManagementPolicy `json:"faultManagementPolicy"`
	DescriptorReference string `json:"descriptor_reference"`
	ResourceReservation string `json:"resource_reservation"`
	RuntimePolicyInfo string `json:"runtime_policy_info"`
	Status Status `json:"status"`
	Notification string `json:"notification"` 
	LifecycleEventHistory []*LifecycleEvent `json:"lifecycle_event_history"`
	AuditLog string `json:"audit_log"`
	CreatedAt string `json:"createdAt"`
	KeyNames []string `json:"keyNames"`
	Name string `json:"name"`
}

// PhysicalNetworkFunctionRecord based on ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type PhysicalNetworkFunctionRecord struct {
	ID string `json:"id"`
	Vendor string `json:"vendor"`
	Version string `json:"version"`
	Description string `json:"description"`
	ProjectID string `json:"projectId"`
	ConnectionPoint []*ConnectionPoint `json:"connection_point"`
	ParentNSID string `json:"parent_ns_id"`
	DescriptorReference string `json:"descriptor_reference"`
	VNFFGR []*VNFForwardingGraphRecord `json:"vnffgr"`
	OamReference string `json:"oam_reference"`
	ConnectedVirtualLink []*VirtualLinkRecord `json:"connected_virtual_link"`
	PNFAddress []string `json:"pnf_address"`
}

type Policy struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
}