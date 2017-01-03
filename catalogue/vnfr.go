package catalogue

import (
	"fmt"
)

// A VirtualNetworkFunctionRecord as described by ETSI GS NFV-MAN 001 V1.1.1
type VirtualNetworkFunctionRecord struct {
	ID                            ID                       `json:"id"`
	HbVersion                     int                      `json:"hb_version"`
	AutoScalePolicies             []*AutoScalePolicy       `json:"auto_scale_policy"`
	ConnectionPoints              []*ConnectionPoint       `json:"connection_point"`
	ProjectID                     string                   `json:"projectId"`
	DeploymentFlavourKey          string                   `json:"deployment_flavour_key"`
	Configurations                *Configuration           `json:"configurations,omitempty"`
	LifecycleEvents               []*LifecycleEvent        `json:"lifecycle_event"`
	LifecycleEventHistory         []*HistoryLifecycleEvent `json:"lifecycle_event_history"`
	Localization                  string                   `json:"localization"`
	MonitoringParameters          []string                 `json:"monitoring_parameter"`
	VDUs                          []*VirtualDeploymentUnit `json:"vdu"`
	Vendor                        string                   `json:"vendor"`
	Version                       string                   `json:"version"`
	VirtualLinks                  []*InternalVirtualLink   `json:"virtual_link"`
	ParentNsID                    string                   `json:"parent_ns_id"`
	DescriptorReference           ID                       `json:"descriptor_reference"`
	VNFMID                        string                   `json:"vnfm_id"`
	ConnectedExternalVirtualLinks []*VirtualLinkRecord     `json:"connected_external_virtual_link"`
	VNFAddresses                  []string                 `json:"vnf_address"`
	Status                        Status                   `json:"status"`
	Notifications                 []string                 `json:"notification"`
	AuditLog                      string                   `json:"audit_log"`
	RuntimePolicyInfos            []string                 `json:"runtime_policy_info"`
	Name                          string                   `json:"name"`
	Type                          string                   `json:"type"`
	Endpoint                      string                   `json:"endpoint"`
	Task                          string                   `json:"task"`
	Requires                      *Configuration           `json:"requires,omitempty"`
	Provides                      *Configuration           `json:"provides,omitempty"`
	CyclicDependency              bool                     `json:"cyclic_dependency"`
	PackageID                     string                   `json:"packageId"`
}

func NewVNFR(
	vnfd *VirtualNetworkFunctionDescriptor,
	flavourKey string,
	vlrs []*VirtualLinkRecord,
	extension map[string]string,
	vimInstances map[ID][]*VIMInstance) (*VirtualNetworkFunctionRecord, error) {

	autoScalePolicies := make([]*AutoScalePolicy, len(vnfd.AutoScalePolicies))
	for _, asp := range vnfd.AutoScalePolicies {
		autoScalePolicies = append(autoScalePolicies, cloneAutoScalePolicy(asp, vnfd))
	}

	configurations := &Configuration{
		Name: vnfd.Name,
		ConfigurationParameters: []*ConfigurationParameter{},
	}

	if vnfd.Configurations != nil {
		configurations.Name = vnfd.Configurations.Name

		for _, confParam := range vnfd.Configurations.ConfigurationParameters {
			configurations.Append(&ConfigurationParameter{
				ID:      GenerateID(),
				ConfKey: confParam.ConfKey,
				Value:   confParam.Value,
			})
		}
	}

	connectionPoints := make([]*ConnectionPoint, len(vnfd.ConnectionPoints))
	for i, connectionPoint := range vnfd.ConnectionPoints {
		connectionPoints[i] = new(ConnectionPoint)
		*connectionPoints[i] = *connectionPoint
	}

	var endpoint string
	if vnfd.Endpoint != "" {
		endpoint = vnfd.Endpoint
	} else {
		endpoint = vnfd.Type
	}

	lifecycleEvents := make([]*LifecycleEvent, len(vnfd.LifecycleEvents))
	for i, lifecycleEvent := range vnfd.LifecycleEvents {
		lceStrings := make([]string, len(lifecycleEvent.LifecycleEvents))
		copy(lceStrings, lifecycleEvent.LifecycleEvents)

		lifecycleEvents[i] = &LifecycleEvent{
			ID:              GenerateID(),
			Event:           lifecycleEvent.Event,
			LifecycleEvents: lceStrings,
		}
	}

	monitoringParameters := make([]string, len(vnfd.MonitoringParameters))
	copy(monitoringParameters, vnfd.MonitoringParameters)

	nsrID := extension["nsr-id"]

	provides := &Configuration{
		Name: "provides",
		ConfigurationParameters: []*ConfigurationParameter{},
	}

	if vnfd.Provides != nil {
		for _, key := range vnfd.Provides {
			provides.Append(&ConfigurationParameter{
				ID:      GenerateID(),
				ConfKey: key,
			})
		}
	}

	requires := &Configuration{
		ID:   GenerateID(),
		Name: "requires",
		ConfigurationParameters: []*ConfigurationParameter{},
	}

	if vnfd.Requires != nil {
		for _, requiresParam := range vnfd.Requires {
			for _, key := range requiresParam.Parameters {
				requires.Append(&ConfigurationParameter{
					ID:      GenerateID(),
					ConfKey: key,
				})
			}
		}
	}

	vdus := make([]*VirtualDeploymentUnit, len(vnfd.VDUs))
	for i, vdu := range vnfd.VDUs {
		for _, vi := range vimInstances[vdu.ID] {
			for _, name := range vdu.VIMInstanceNames {
				if name == vi.Name {
					if !vi.HasFlavour(flavourKey) {
						return nil, fmt.Errorf("no key %s found in vim instance: %v", flavourKey, vi)
					}
				}
			}
		}

		vdus[i] = makeVDUFromParent(vdu)
	}

	links := make([]*InternalVirtualLink, len(vnfd.VirtualLinks))
	for i, oldIVL := range vnfd.VirtualLinks {
		links[i] = cloneInternalVirtualLink(oldIVL, vlrs)
	}

	return &VirtualNetworkFunctionRecord{
		ID:   GenerateID(),
		Name: vnfd.Name,

		AutoScalePolicies:     autoScalePolicies,
		Configurations:        configurations,
		ConnectionPoints:      connectionPoints,
		CyclicDependency:      vnfd.CyclicDependency,
		DeploymentFlavourKey:  flavourKey,
		DescriptorReference:   vnfd.ID,
		Endpoint:              endpoint,
		LifecycleEventHistory: []*HistoryLifecycleEvent{},
		LifecycleEvents:       lifecycleEvents,
		MonitoringParameters:  monitoringParameters,
		PackageID:             vnfd.VNFPackageLocation,
		ParentNsID:            nsrID,
		Provides:              provides,
		Requires:              requires,
		Status:                StatusNull,
		Type:                  vnfd.Type,
		Vendor:                vnfd.Vendor,
		Version:               vnfd.Version,
		VNFAddresses:          []string{},
	}, nil

	// TODO mange the VirtualLinks and links...
}

// FindComponentInstance searches an instance of a given VNFComponent inside the
// VirtualDeploymentUnit of the current VirtualNetworkFunctionRecord.
func (vnfr *VirtualNetworkFunctionRecord) FindComponentInstance(component *VNFComponent) *VNFCInstance {
	for _, vdu := range vnfr.VDUs {
		for _, vnfcInstance := range vdu.VNFCInstances {
			if vnfcInstance.VNFComponent.ID == component.ID {
				return vnfcInstance
			}
		}
	}

	return nil
}

func cloneAutoScalePolicy(asp *AutoScalePolicy, vnfd *VirtualNetworkFunctionDescriptor) *AutoScalePolicy {
	// copy all in bulk, and then deep clone the pointers
	newAsp := new(AutoScalePolicy)
	*newAsp = *asp

	newAsp.Actions = make([]*ScalingAction, len(asp.Actions))
	for i, action := range asp.Actions {
		target := action.Target
		if target == "" {
			target = vnfd.Type
		}

		newAsp.Actions[i] = &ScalingAction{
			ID:     GenerateID(),
			Target: target,
			Type:   action.Type,
			Value:  action.Value,
		}
	}

	newAsp.Alarms = make([]*ScalingAlarm, len(asp.Alarms))
	for i, alarm := range asp.Alarms {
		newAsp.Alarms[i] = &ScalingAlarm{
			ID:                 GenerateID(),
			ComparisonOperator: alarm.ComparisonOperator,
			Metric:             alarm.Metric,
			Statistic:          alarm.Statistic,
			Threshold:          alarm.Threshold,
			Weight:             alarm.Weight,
		}
	}

	return newAsp
}

func cloneInternalVirtualLink(oldIVL *InternalVirtualLink, vlrs []*VirtualLinkRecord) *InternalVirtualLink {
	extID := ""
	name := oldIVL.Name

	for _, vlr := range vlrs {
		if vlr.Name == name {
			extID = vlr.ExtID
		}
	}

	cpReferences := make([]string, len(oldIVL.ConnectionPointsReferences))
	copy(cpReferences, oldIVL.ConnectionPointsReferences)

	qos := make([]string, len(oldIVL.QoS))
	copy(qos, oldIVL.QoS)

	testAccess := make([]string, len(oldIVL.TestAccess))
	copy(testAccess, oldIVL.TestAccess)

	return &InternalVirtualLink{
		VirtualLink: VirtualLink{
			ID:               GenerateID(),
			Name:             name,
			ConnectivityType: oldIVL.ConnectivityType,
			ExtID:            extID,
			LeafRequirement:  oldIVL.LeafRequirement,
			QoS:              qos,
			RootRequirement:  oldIVL.RootRequirement,
			TestAccess:       testAccess,
		},

		ConnectionPointsReferences: cpReferences,
	}
}

func cloneVRFaultManagementPolicy(oldVRFMP *VRFaultManagementPolicy) *VRFaultManagementPolicy {
	newVRFMP := new(VRFaultManagementPolicy)
	*newVRFMP = *oldVRFMP

	newVRFMP.Criteria = make([]*Criteria, len(oldVRFMP.Criteria))
	for i, criteria := range oldVRFMP.Criteria {
		newVRFMP.Criteria[i] = new(Criteria)
		*newVRFMP.Criteria[i] = *criteria
	}

	return newVRFMP
}

func makeVDUFromParent(parentVDU *VirtualDeploymentUnit) *VirtualDeploymentUnit {
	// copy all of the struct at once, and then deep clone the pointer/list parts
	newVDU := new(VirtualDeploymentUnit)
	*newVDU = *parentVDU

	// reset the ID of the new VDU
	newVDU.ID = GenerateID()
	newVDU.ParentVDU = parentVDU.ID

	newVDU.VNFCs = make([]*VNFComponent, len(parentVDU.VNFCs))

	for i, component := range parentVDU.VNFCs {
		connectionPoints := make([]*VNFDConnectionPoint, len(parentVDU.VNFCs))
		for j, connectionPoint := range component.ConnectionPoints {
			connectionPoints[j] = &VNFDConnectionPoint{
				ConnectionPoint: ConnectionPoint{
					ID:   GenerateID(),
					Type: connectionPoint.Type,
				},

				FloatingIP:           connectionPoint.FloatingIP,
				VirtualLinkReference: connectionPoint.VirtualLinkReference,
			}
		}

		newVDU.VNFCs[i] = &VNFComponent{
			ID:               GenerateID(),
			ConnectionPoints: connectionPoints,
		}
	}

	newVDU.VNFCInstances = make([]*VNFCInstance, len(parentVDU.VNFCInstances))

	newVDU.LifecycleEvents = make([]*LifecycleEvent, len(parentVDU.LifecycleEvents))

	for i, lifecycleEvent := range parentVDU.LifecycleEvents {
		lifecycleEvents := make([]string, len(lifecycleEvent.LifecycleEvents))
		copy(lifecycleEvents, lifecycleEvent.LifecycleEvents)

		newVDU.LifecycleEvents[i] = &LifecycleEvent{
			ID:              GenerateID(),
			Event:           lifecycleEvent.Event,
			LifecycleEvents: lifecycleEvents,
		}
	}

	newVDU.MonitoringParameters = make([]string, len(parentVDU.MonitoringParameters))
	copy(newVDU.MonitoringParameters, parentVDU.MonitoringParameters)

	newVDU.FaultManagementPolicies = make([]*VRFaultManagementPolicy, len(parentVDU.FaultManagementPolicies))

	if parentVDU.FaultManagementPolicies != nil {
		for i, vrfmp := range parentVDU.FaultManagementPolicies {
			newVDU.FaultManagementPolicies[i] = cloneVRFaultManagementPolicy(vrfmp)
		}
	}

	newVDU.VMImages = make([]string, len(parentVDU.VMImages))
	copy(newVDU.VMImages, parentVDU.VMImages)

	newVDU.VIMInstanceNames = make([]string, len(parentVDU.VIMInstanceNames))
	copy(newVDU.VIMInstanceNames, parentVDU.VIMInstanceNames)

	if parentVDU.HighAvailability != nil {
		newVDU.HighAvailability = new(HighAvailability)
		*newVDU.HighAvailability = *parentVDU.HighAvailability
	}

	return newVDU
}
