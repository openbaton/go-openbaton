package openbaton

type FaultManagementAction string

const (
	FaultRestart              FaultManagementAction = "RESTART"
	FaultReinstantiateService FaultManagementAction = "REINSTANTIATE_SERVICE"
	FaultHeal                 FaultManagementAction = "HEAL"
	FaultReinstantiate        FaultManagementAction = "REINSTANTIATE"
	FaultSwitchToStandby      FaultManagementAction = "SWITCH_TO_STANDBY"
	FaultSwitchToActive       FaultManagementAction = "SWITCH_TO_ACTIVE"
)

type VRFaultManagementPolicy struct {
	Action FaultManagementAction `json:"action"`
}