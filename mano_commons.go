package openbaton

type AutoScalePolicy struct {
	ID                 string          `json:"id"`
	Version            int             `json:"version"`
	Name               string          `json:"name"`
	Threshold          float64         `json:"threshold"`
	ComparisonOperator string          `json:"comparisonOperator"`
	Period             int             `json:"period"`
	Cooldown           int             `json:"cooldown"`
	Mode               ScalingMode     `json:"mode"`
	Type               ScalingType     `json:"type"`
	Alarms             []*ScalingAlarm `json:"alarms"`
	Actions            []ScalingAction `json:"actions"`
}

type ConnectionPoint struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
	Type    string `json:"type"`
}

// ConstituentVNF as described in ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type ConstituentVNF struct {
	ID                    string          `json:"id"`
	VnfReference          string          `json:"vnf_reference"`
	VnfFlavourIDReference string          `json:"vnf_flavour_id_reference"`
	RedundancyModel       RedundancyModel `json:"redundancy_modelid"`
	Affinity              string          `json:"affinity"`
	Capability            string          `json:"capability"`
	NumberOfInstances     int             `json:"number_of_instancesid"`
	Version               int             `json:"version"`
}

type Event string

const (
	EventGranted  Event = "GRANTED"
	EventAllocate Event = "ALLOCATE"
	EventScale    Event = "SCALE"
	EventRelease  Event = "RELEASE"
	EventError    Event = "ERROR"

	EventInstantiate     Event = "INSTANTIATE"
	EventTerminate       Event = "TERMINATE"
	EventConfigure       Event = "CONFIGURE"
	EventStart           Event = "START"
	EventStop            Event = "STOP"
	EventHeal            Event = "HEAL"
	EventScaleOut        Event = "SCALE_OUT"
	EventScaleIn         Event = "SCALE_IN"
	EventScaleUp         Event = "SCALE_UP"
	EventScaleDown       Event = "SCALE_DOWN"
	EventUpdate          Event = "UPDATE"
	EventUpdateRollback  Event = "UPDATE_ROLLBACK"
	EventUpgrade         Event = "UPGRADE"
	EventUpgradeRollback Event = "UPGRADE_ROLLBACK"
	EventReset           Event = "RESET"
)

type FaultManagementAction string

const (
	FaultRestart              FaultManagementAction = "RESTART"
	FaultReinstantiateService FaultManagementAction = "REINSTANTIATE_SERVICE"
	FaultHeal                 FaultManagementAction = "HEAL"
	FaultReinstantiate        FaultManagementAction = "REINSTANTIATE"
	FaultSwitchToStandby      FaultManagementAction = "SWITCH_TO_STANDBY"
	FaultSwitchToActive       FaultManagementAction = "SWITCH_TO_ACTIVE"
)

type HighAvailability struct {
	ID               string          `json:"id"`
	Version          int             `json:"version"`
	ResiliencyLevel  ResiliencyLevel `json:"resiliencyLevel"`
	GeoRedundancy    bool            `json:"geoRedundancy"`
	RedundancyScheme string          `json:"redundancyScheme"`
}

type Ip struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
	NetName string `json:"netname"`
	Ip      string `json:"ip"`
}

// A LifecycleEvent as specified in ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type LifecycleEvent struct {
	ID              string   `json:"id"`
	Version         int      `json:"version"`
	Event           Event    `json:"event"`
	LifecycleEvents []string `json:"lifecycle_events"`
}

type RedundancyModel string

const (
	RedundancyActive  RedundancyModel = "ACTIVE"
	RedundancyStandby RedundancyModel = "STANDBY"
)

type ResiliencyLevel string

const (
	ResiliencyActiveStandbyStateless ResiliencyLevel = "ACTIVE_STANDBY_STATELESS"
	ResiliencyActiveStandbyStateful  ResiliencyLevel = "ACTIVE_STANDBY_STATEFUL"
)

type ScalingAlarm struct {
	ID                 string  `json:"id"`
	Version            int     `json:"version"`
	Metric             string  `json:"metric"`
	Statistic          string  `json:"statistic"`
	ComparisonOperator string  `json:"comparisonOperator"`
	Threshold          float64 `json:"threshold"`
	Weight             float64 `json:"weight"`
}

type ScalingAction struct {
	ID      string            `json:"id"`
	Version int               `json:"version"`
	Type    ScalingActionType `json:"type"`
	Value   string            `json:"value"`
	Target  string            `json:"target"`
}

type ScalingActionType string

const (
	ScaleOut          ScalingActionType = "SCALE_OUT"
	ScaleOutTo        ScalingActionType = "SCALE_OUT_TO"
	ScaleOutToFlavour ScalingActionType = "SCALE_OUT_TO_FLAVOUR"
	ScaleIn           ScalingActionType = "SCALE_IN"
	ScaleInTo         ScalingActionType = "SCALE_IN_TO"
	ScaleInToFlavour  ScalingActionType = "SCALE_IN_TO_FLAVOUR"
)

type ScalingMode string

const (
	ScaleModeReactive   ScalingMode = "REACTIVE"
	ScaleModeProactive  ScalingMode = "PROACTIVE"
	ScaleModePredictive ScalingMode = "PREDICTIVE"
)

type ScalingType string

const (
	ScaleTypeSingle   ScalingType = "SINGLE"
	ScaleTypeVoted    ScalingType = "VOTED"
	ScaleTypeWeighted ScalingType = "WEIGHTED"
)

type Security struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
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

// VirtualLink (abstract) based on ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
// The VLD describes the basic topology of the connectivity (e.g. E-LAN, E-Line, E-Tree) between one
// or more VNFs connected to this VL and other required parameters (e.g. bandwidth and QoS class).
// The VLD connection parameters are expected to have similar attributes to those used on the ports
// on VNFs in ETSI GS NFV-SWA 001 [i.8]. Therefore a set of VLs in a Network Service can be mapped
// to a Network Connectivity Topology (NCT) as defined in ETSI GS NFV-SWA 001 [i.8].
type VirtualLink struct {
	ID               string   `json:"id"`
	HbVersion        int      `json:"hb_version"`
	ExtID            string   `json:"extId"`
	RootRequirement  string   `json:"root_requirement"`
	LeafRequirement  string   `json:"leaf_requirement"`
	Qos              []string `json:"qos"`
	TestAccess       []string `json:"test_access"`
	ConnectivityType []string `json:"connectivity_type"`
	Name             string   `json:"name"`
}

// VirtualLinkDescriptor as described in ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type VirtualLinkDescriptor struct {
	VirtualLink
	ProjectID         string    `json:"projectId"`
	Vendor            string    `json:"vendor"`
	DescriptorVersion string    `json:"descriptor_version"`
	NumberOfEndpoints int       `json:"number_of_endpoints"`
	connection        []string  `json:"connection"`
	vld_security      *Security `json:"vld_security"`
	name              string
}

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

type VNFCInstance struct {
	VimID        string        `json:"vim_id"`
	VcID         string        `json:"vc_id"`
	Hostname     string        `json:"hostname"`
	State        string        `json:"state"`
	VnfComponent *VNFComponent `json:"vnfComponent"`
	FloatingIps  []*Ip         `json:"floatingIps"`
	Ips          []*Ip         `json:"ips"`
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

// A Virtual Network Function Record as described by ETSI GS NFV-MAN 001 V1.1.1
type VNFRecord struct {
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

type VRFaultManagementPolicy struct {
	Action FaultManagementAction `json:"action"`
}
