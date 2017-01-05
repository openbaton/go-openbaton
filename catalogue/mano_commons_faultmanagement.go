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

