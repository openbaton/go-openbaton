package catalogue

type AutoScalePolicy struct {
	ID                 ID               `json:"id"`
	Version            int              `json:"version"`
	Name               string           `json:"name"`
	Threshold          float64          `json:"threshold"`
	ComparisonOperator string           `json:"comparisonOperator"`
	Period             int              `json:"period"`
	Cooldown           int              `json:"cooldown"`
	Mode               ScalingMode      `json:"mode,omitempty"`
	Type               ScalingType      `json:"type,omitempty"`
	Alarms             []*ScalingAlarm  `json:"alarms,omitempty"`
	Actions            []*ScalingAction `json:"actions,omitempty"`
}

type ConnectionPoint struct {
	ID      ID     `json:"id"`
	Version int    `json:"version"`
	Type    string `json:"type"`
}

// ConstituentVDU as described in ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type ConstituentVDU struct {
	ID                ID     `json:"id"`
	Version           int    `json:"version"`
	VDUReference      string `json:"vdu_reference"`
	NumberOfInstances int    `json:"number_of_instances"`
	ConstituentVNFC   string `json:"constituent_vnfc"`
}

// ConstituentVNF as described in ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type ConstituentVNF struct {
	ID                    ID              `json:"id"`
	VnfReference          string          `json:"vnf_reference"`
	VnfFlavourIDReference string          `json:"vnf_flavour_id_reference"`
	RedundancyModel       RedundancyModel `json:"redundancy_modelid"`
	Affinity              string          `json:"affinity"`
	Capability            string          `json:"capability"`
	NumberOfInstances     int             `json:"number_of_instancesid"`
	Version               int             `json:"version"`
}

// DeploymentFlavour as described in ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type DeploymentFlavour struct {
	ID         ID     `json:"id"`
	Version    int    `json:"version"`
	FlavourKey string `json:"flavour_key"`
	ExtID      string `json:"extId"`
	RAM        int    `json:"ram"`
	Disk       int    `json:"disk"`
	VCPUs      int    `json:"vcpus"`
}

type Event string

const (
	EventGranted  = Event("GRANTED")
	EventAllocate = Event("ALLOCATE")
	EventScale    = Event("SCALE")
	EventRelease  = Event("RELEASE")
	EventError    = Event("ERROR")

	EventInstantiate     = Event("INSTANTIATE")
	EventTerminate       = Event("TERMINATE")
	EventConfigure       = Event("CONFIGURE")
	EventStart           = Event("START")
	EventStop            = Event("STOP")
	EventHeal            = Event("HEAL")
	EventScaleOut        = Event("SCALE_OUT")
	EventScaleIn         = Event("SCALE_IN")
	EventScaleUp         = Event("SCALE_UP")
	EventScaleDown       = Event("SCALE_DOWN")
	EventUpdate          = Event("UPDATE")
	EventUpdateRollback  = Event("UPDATE_ROLLBACK")
	EventUpgrade         = Event("UPGRADE")
	EventUpgradeRollback = Event("UPGRADE_ROLLBACK")
	EventReset           = Event("RESET")
)

type HighAvailability struct {
	ID               ID              `json:"id"`
	Version          int             `json:"version"`
	ResiliencyLevel  ResiliencyLevel `json:"resiliencyLevel"`
	GeoRedundancy    bool            `json:"geoRedundancy"`
	RedundancyScheme string          `json:"redundancyScheme"`
}

type Ip struct {
	ID      ID     `json:"id"`
	Version int    `json:"version"`
	NetName string `json:"netname"`
	Ip      string `json:"ip"`
}

// LifecycleEvent as specified in ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type LifecycleEvent struct {
	ID              ID       `json:"id"`
	Version         int      `json:"version"`
	Event           Event    `json:"event"`
	LifecycleEvents []string `json:"lifecycle_events"`
}

// NetworkServiceDeploymentFlavour as specified in ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type NetworkServiceDeploymentFlavour struct {
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

type RedundancyModel string

const (
	RedundancyActive  = RedundancyModel("ACTIVE")
	RedundancyStandby = RedundancyModel("STANDBY")
)

type ResiliencyLevel string

const (
	ResiliencyActiveStandbyStateless = ResiliencyLevel("ACTIVE_STANDBY_STATELESS")
	ResiliencyActiveStandbyStateful  = ResiliencyLevel("ACTIVE_STANDBY_STATEFUL")
)

type ScalingAction struct {
	ID      ID                `json:"id"`
	Version int               `json:"version"`
	Type    ScalingActionType `json:"type"`
	Value   string            `json:"value,omitempty"`
	Target  string            `json:"target,omitempty"`
}

type ScalingActionType string

const (
	ScaleOut          = ScalingActionType("SCALE_OUT")
	ScaleOutTo        = ScalingActionType("SCALE_OUT_TO")
	ScaleOutToFlavour = ScalingActionType("SCALE_OUT_TO_FLAVOUR")
	ScaleIn           = ScalingActionType("SCALE_IN")
	ScaleInTo         = ScalingActionType("SCALE_IN_TO")
	ScaleInToFlavour  = ScalingActionType("SCALE_IN_TO_FLAVOUR")
)

type ScalingAlarm struct {
	ID                 ID      `json:"id"`
	Version            int     `json:"version"`
	Metric             string  `json:"metric"`
	Statistic          string  `json:"statistic"`
	ComparisonOperator string  `json:"comparisonOperator"`
	Threshold          float64 `json:"threshold"`
	Weight             float64 `json:"weight"`
}

type ScalingMode string

const (
	ScaleModeReactive   = ScalingMode("REACTIVE")
	ScaleModeProactive  = ScalingMode("PROACTIVE")
	ScaleModePredictive = ScalingMode("PREDICTIVE")
)

type ScalingType string

const (
	ScaleTypeSingle   = ScalingType("SINGLE")
	ScaleTypeVoted    = ScalingType("VOTED")
	ScaleTypeWeighted = ScalingType("WEIGHTED")
)

type Security struct {
	ID      ID  `json:"id"`
	Version int `json:"version"`
}

// VirtualLink (abstract) based on ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
// The VLD describes the basic topology of the connectivity (e.g. E-LAN, E-Line, E-Tree) between one
// or more VNFs connected to this VL and other required parameters (e.g. bandwidth and QoS class).
// The VLD connection parameters are expected to have similar attributes to those used on the ports
// on VNFs in ETSI GS NFV-SWA 001 [i.8]. Therefore a set of VLs in a Network Service can be mapped
// to a Network Connectivity Topology (NCT) as defined in ETSI GS NFV-SWA 001 [i.8].
type VirtualLink struct {
	ID               ID       `json:"id"`
	HbVersion        int      `json:"hb_version"`
	ExtID            string   `json:"extId"`
	RootRequirement  string   `json:"root_requirement"`
	LeafRequirement  string   `json:"leaf_requirement"`
	QoS              []string `json:"qos"`
	TestAccess       []string `json:"test_access"`
	ConnectivityType []string `json:"connectivity_type"`
	Name             string   `json:"name"`
}

type VNFDeploymentFlavour struct {
	DeploymentFlavour
	DfConstraint   []string          `json:"df_constraint"`
	ConstituentVDU []*ConstituentVDU `json:"constituent_vdu"`
}
