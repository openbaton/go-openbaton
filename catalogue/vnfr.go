package catalogue

// A VirtualNetworkFunctionRecord as described by ETSI GS NFV-MAN 001 V1.1.1
type VirtualNetworkFunctionRecord struct {
	ID                           ID                   `json:"id"`
	HbVersion                    int                      `json:"hb_version"`
	AutoScalePolicies              []*AutoScalePolicy       `json:"auto_scale_policy"`
	ConnectionPoints              []*ConnectionPoint       `json:"connection_point"`
	ProjectID                    string                   `json:"projectId"`
	DeploymentFlavourKey         string                   `json:"deployment_flavour_key"`
	Configurations               *Configuration           `json:"configurations,omitempty"`
	LifecycleEvents               []*LifecycleEvent        `json:"lifecycle_event"`
	LifecycleEventHistory        []*HistoryLifecycleEvent `json:"lifecycle_event_history"`
	Localization                 string                   `json:"localization"`
	MonitoringParameters          []string                 `json:"monitoring_parameter"`
	VDUs                         []*VirtualDeploymentUnit `json:"vdu"`
	Vendor                       string                   `json:"vendor"`
	Version                      string                   `json:"version"`
	VirtualLinks                  []InternalVirtualLink    `json:"virtual_link"`
	ParentNsID                   string                   `json:"parent_ns_id"`
	DescriptorReference          string                   `json:"descriptor_reference"`
	VNFMID                       string                   `json:"vnfm_id"`
	ConnectedExternalVirtualLinks []VirtualLinkRecord      `json:"connected_external_virtual_link"`
	VnfAddresses                   []string                 `json:"vnf_address"`
	Status                       Status                   `json:"status"`
	Notifications                 []string                 `json:"notification"`
	AuditLog                     string                   `json:"audit_log"`
	RuntimePolicyInfos            []string                 `json:"runtime_policy_info"`
	Name                         string                   `json:"name"`
	Type                         string                   `json:"type"`
	Endpoint                     string                   `json:"endpoint"`
	Task                         string                   `json:"task"`
	Requires                     *Configuration           `json:"requires,omitempty"`
	Provides                     *Configuration           `json:"provides,omitempty"`
	CyclicDependency             bool                     `json:"cyclic_dependency"`
	PackageID                    string                   `json:"packageId"`
}

func NewVNFR(
      vnfd *VirtualNetworkFunctionDescriptor,
      flavourKey string,
	  vlr []*VirtualLinkRecord,
	  extension map[string]string,
	  vimInstances    map[string][]*catalogue.VIMInstance) (*VirtualNetworkFunctionRecord, error) {

    autoScalePolicies := make([]*AutoScalePolicy, len(vnfd.AutoScalePolicies))
    for _, asp := vnfd.AutoScalePolicies {
      autoScalePolicies = append(autoScalePolicies, cloneAutoScalePolicy(asp, vnfd))
    }

    configurations := &Configuration{
        Name: vnfd.Name,
        ConfigurationParameters: []*ConfigurationParameter{},
    }

    if vnfd.Configurations != nil {
        configurations.Name = vnfd.Configurations.Name

        for _, confParam := range vnfd.Configurations.ConfigurationParameters {
            configurations.ConfigurationParameter.Append(&ConfigurationParameter{
                ConfKey: confParam.ConfKey,
                Value: confParam.Value,
            })
        }
    }

    var endpoint string
    if vnfd.Endpoint != "" {
      endpoint = vnfd.Endpoint
    } else {
      endpoint = vnfd.Type
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
        requires.Append(&ConfigurationParameter{
          ConfKey: key,
        })
      }
    }

    requires := &Configuration{
      Name: "requires",
      ConfigurationParameters: []*ConfigurationParameter{},
    }

    if vnfd.Requires != nil {
      for _, requiresParam := range vnfd.Requires {
        for _, key := range requiresParam.Parameters {
          requires.Append(&ConfigurationParameter{ConfKey: key})
        }
      }
    }

    vdus := make([]*VirtualDeploymentUnit, len(vnfd.VDUs))
    for _, vdu := vnfd.VDUs {
      vdus = append(vdus, cloneVDU(vdu))
    }

    vnfr := &VirtualNetworkFunctionRecord{
        Name: vnfd.Name,

        AutoScalePolicies: autoScalePolicies,
        Configurations: configurations,
        CyclicDependency: vnfd.CyclicDependency,
        Endpoint: endpoint,
        LifecycleEventHistory: []*HistoryLifecycleEvent{},
        MonitoringParameters: monitoringParameters,
        PackageID: vnfd.VNFPackageLocation,
        ParentNsID: nsrID,
        Provides: provides,
        Requires: requires,
        Type: vnfd.Type,
        Vendor: vnfd.Vendor,
    }

    // TODO mange the VirtualLinks and links...

    virtualNetworkFunctionRecord.setVdu(new HashSet<VirtualDeploymentUnit>());
    for (VirtualDeploymentUnit virtualDeploymentUnit : vnfd.getVdu()) {
      VirtualDeploymentUnit vdu_new = new VirtualDeploymentUnit();
      vdu_new.setParent_vdu(virtualDeploymentUnit.getId());

      HashSet<VNFComponent> vnfComponents = new HashSet<>();
      for (VNFComponent component : virtualDeploymentUnit.getVnfc()) {
        VNFComponent component_new = new VNFComponent();
        HashSet<VNFDConnectionPoint> connectionPoints = new HashSet<>();
        for (VNFDConnectionPoint connectionPoint : component.getConnection_point()) {
          VNFDConnectionPoint connectionPoint_new = new VNFDConnectionPoint();
          connectionPoint_new.setVirtual_link_reference(
              connectionPoint.getVirtual_link_reference());
          connectionPoint_new.setType(connectionPoint.getType());
          connectionPoint_new.setFloatingIp(connectionPoint.getFloatingIp());
          connectionPoints.add(connectionPoint_new);
        }
        component_new.setConnection_point(connectionPoints);
        vnfComponents.add(component_new);
      }
      vdu_new.setVnfc(vnfComponents);
      vdu_new.setVnfc_instance(new HashSet<VNFCInstance>());
      HashSet<LifecycleEvent> lifecycleEvents = new HashSet<>();
      for (LifecycleEvent lifecycleEvent : virtualDeploymentUnit.getLifecycle_event()) {
        LifecycleEvent lifecycleEvent_new = new LifecycleEvent();
        lifecycleEvent_new.setEvent(lifecycleEvent.getEvent());
        lifecycleEvent_new.setLifecycle_events(lifecycleEvent.getLifecycle_events());
        lifecycleEvents.add(lifecycleEvent_new);
      }
      vdu_new.setLifecycle_event(lifecycleEvents);
      vdu_new.setVimInstanceName(virtualDeploymentUnit.getVimInstanceName());
      vdu_new.setHostname(virtualDeploymentUnit.getHostname());
      vdu_new.setHigh_availability(virtualDeploymentUnit.getHigh_availability());
      vdu_new.setComputation_requirement(virtualDeploymentUnit.getComputation_requirement());
      vdu_new.setScale_in_out(virtualDeploymentUnit.getScale_in_out());
      HashSet<String> monitoringParameters = new HashSet<>();
      monitoringParameters.addAll(virtualDeploymentUnit.getMonitoring_parameter());
      vdu_new.setMonitoring_parameter(monitoringParameters);
      vdu_new.setVdu_constraint(virtualDeploymentUnit.getVdu_constraint());

      //Set Faultmanagement policies
      Set<VRFaultManagementPolicy> vrFaultManagementPolicies = new HashSet<>();
      if (virtualDeploymentUnit.getFault_management_policy() != null) {
        log.debug(
            "Adding the fault management policies: "
                + virtualDeploymentUnit.getFault_management_policy());
        for (VRFaultManagementPolicy vrfmp : virtualDeploymentUnit.getFault_management_policy()) {
          vrFaultManagementPolicies.add(vrfmp);
        }
      }
      vdu_new.setFault_management_policy(vrFaultManagementPolicies);
      //Set Faultmanagement policies end

      HashSet<String> vmImages = new HashSet<>();
      vmImages.addAll(virtualDeploymentUnit.getVm_image());
      vdu_new.setVm_image(vmImages);

      vdu_new.setVirtual_network_bandwidth_resource(
          virtualDeploymentUnit.getVirtual_network_bandwidth_resource());
      vdu_new.setVirtual_memory_resource_element(
          virtualDeploymentUnit.getVirtual_memory_resource_element());
      virtualNetworkFunctionRecord.getVdu().add(vdu_new);
    }
    virtualNetworkFunctionRecord.setVersion(vnfd.getVersion());
    virtualNetworkFunctionRecord.setConnection_point(new HashSet<ConnectionPoint>());
    virtualNetworkFunctionRecord.getConnection_point().addAll(vnfd.getConnection_point());

    // TODO find a way to choose between deployment flavors and create the new one
    virtualNetworkFunctionRecord.setDeployment_flavour_key(flavourKey);
    for (VirtualDeploymentUnit virtualDeploymentUnit : vnfd.getVdu()) {
      for (VimInstance vi : vimInstances.get(virtualDeploymentUnit.getId())) {
        for (String name : virtualDeploymentUnit.getVimInstanceName()) {
          if (name.equals(vi.getName())) {
            if (!existsDeploymentFlavor(
                virtualNetworkFunctionRecord.getDeployment_flavour_key(), vi)) {
              throw new BadFormatException(
                  "no key "
                      + virtualNetworkFunctionRecord.getDeployment_flavour_key()
                      + " found in vim instance: "
                      + vi);
            }
          }
        }
      }
    }

    virtualNetworkFunctionRecord.setDescriptor_reference(vnfd.getId());
    virtualNetworkFunctionRecord.setLifecycle_event(new LinkedHashSet<LifecycleEvent>());
    HashSet<LifecycleEvent> lifecycleEvents = new HashSet<>();
    for (LifecycleEvent lifecycleEvent : vnfd.getLifecycle_event()) {
      LifecycleEvent lifecycleEvent_new = new LifecycleEvent();
      lifecycleEvent_new.setEvent(lifecycleEvent.getEvent());
      lifecycleEvent_new.setLifecycle_events(new ArrayList<String>());
      for (String event : lifecycleEvent.getLifecycle_events()) {
        lifecycleEvent_new.getLifecycle_events().add(event);
      }
      log.debug(
          "Found SCRIPTS for EVENT "
              + lifecycleEvent_new.getEvent()
              + ": "
              + lifecycleEvent_new.getLifecycle_events().size());
      lifecycleEvents.add(lifecycleEvent_new);
    }
    virtualNetworkFunctionRecord.setLifecycle_event(lifecycleEvents);
    virtualNetworkFunctionRecord.setVirtual_link(new HashSet<InternalVirtualLink>());
    HashSet<InternalVirtualLink> internalVirtualLinks = new HashSet<>();
    for (InternalVirtualLink internalVirtualLink : vnfd.getVirtual_link()) {
      InternalVirtualLink internalVirtualLink_new = new InternalVirtualLink();
      internalVirtualLink_new.setName(internalVirtualLink.getName());

      for (VirtualLinkRecord virtualLinkRecord : vlr) {
        if (virtualLinkRecord.getName().equals(internalVirtualLink_new.getName())) {
          internalVirtualLink_new.setExtId(virtualLinkRecord.getExtId());
        }
      }

      internalVirtualLink_new.setLeaf_requirement(internalVirtualLink.getLeaf_requirement());
      internalVirtualLink_new.setRoot_requirement(internalVirtualLink.getRoot_requirement());
      internalVirtualLink_new.setConnection_points_references(new HashSet<String>());
      for (String conn : internalVirtualLink.getConnection_points_references()) {
        internalVirtualLink_new.getConnection_points_references().add(conn);
      }
      internalVirtualLink_new.setQos(new HashSet<String>());
      for (String qos : internalVirtualLink.getQos()) {
        internalVirtualLink_new.getQos().add(qos);
      }
      internalVirtualLink_new.setTest_access(new HashSet<String>());
      for (String test : internalVirtualLink.getTest_access()) {
        internalVirtualLink_new.getTest_access().add(test);
      }
      internalVirtualLink_new.setConnectivity_type(internalVirtualLink.getConnectivity_type());
      internalVirtualLinks.add(internalVirtualLink_new);
    }
    virtualNetworkFunctionRecord.getVirtual_link().addAll(internalVirtualLinks);

    virtualNetworkFunctionRecord.setVnf_address(new HashSet<String>());
    virtualNetworkFunctionRecord.setStatus(Status.NULL);

	for (InternalVirtualLink internalVirtualLink :
		virtualNetworkFunctionRecord.getVirtual_link()) {
	for (VirtualLinkRecord virtualLinkRecord : virtualLinkRecords) {
		if (internalVirtualLink.getName().equals(virtualLinkRecord.getName())) {
		internalVirtualLink.setExtId(virtualLinkRecord.getExtId());
		internalVirtualLink.setConnectivity_type(virtualLinkRecord.getConnectivity_type());
		}
	}
	}
	log.debug("Created VirtualNetworkFunctionRecordShort: " + virtualNetworkFunctionRecord);
	return virtualNetworkFunctionRecord
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
  newAsp := *asp

  newAsp.Actions = make([]*ScalingAction, len(asp.Actions))
  for _, action := range asp.Actions {
    target := action.Target
    if target == "" {
      target = vnfd.Type
    }

    newAsp.Actions = append(newAsp.Actions, &ScalingAction{
      Target: target,
      Type: action.Type,
      Value: action.Value,
    })
  }
  
  newAsp.Alarms = make([]*ScalingAlarm, len(asp.Alarms))
  for _, alarm := range asp.Alarms {
    newAsp.Alarms = append(newAsp.Alarms, &ScalingAlarm{
      ComparisonOperator: alarm.ComparisonOperator,
      Metric: alarm.Metric,
      Statistic: alarm.Statistic,
      Threshold: alarm.Threshold,
      Weight: alarm.Weight,
    })
  }

  return &newAsp
}

func init(vdu *VirtualDeploymentUnit) *VirtualDeploymentUnit {
  newVDU := *vdu
  
  vdu_new.setParent_vdu(virtualDeploymentUnit.getId());

  HashSet<VNFComponent> vnfComponents = new HashSet<>();
  for (VNFComponent component : virtualDeploymentUnit.getVnfc()) {
    VNFComponent component_new = new VNFComponent();
    HashSet<VNFDConnectionPoint> connectionPoints = new HashSet<>();
    for (VNFDConnectionPoint connectionPoint : component.getConnection_point()) {
      VNFDConnectionPoint connectionPoint_new = new VNFDConnectionPoint();
      connectionPoint_new.setVirtual_link_reference(
          connectionPoint.getVirtual_link_reference());
      connectionPoint_new.setType(connectionPoint.getType());
      connectionPoint_new.setFloatingIp(connectionPoint.getFloatingIp());
      connectionPoints.add(connectionPoint_new);
    }
    component_new.setConnection_point(connectionPoints);
    vnfComponents.add(component_new);
  }
  vdu_new.setVnfc(vnfComponents);
  vdu_new.setVnfc_instance(new HashSet<VNFCInstance>());
  HashSet<LifecycleEvent> lifecycleEvents = new HashSet<>();
  for (LifecycleEvent lifecycleEvent : virtualDeploymentUnit.getLifecycle_event()) {
    LifecycleEvent lifecycleEvent_new = new LifecycleEvent();
    lifecycleEvent_new.setEvent(lifecycleEvent.getEvent());
    lifecycleEvent_new.setLifecycle_events(lifecycleEvent.getLifecycle_events());
    lifecycleEvents.add(lifecycleEvent_new);
  }
  vdu_new.setLifecycle_event(lifecycleEvents);
  vdu_new.setVimInstanceName(virtualDeploymentUnit.getVimInstanceName());
  vdu_new.setHostname(virtualDeploymentUnit.getHostname());
  vdu_new.setHigh_availability(virtualDeploymentUnit.getHigh_availability());
  vdu_new.setComputation_requirement(virtualDeploymentUnit.getComputation_requirement());
  vdu_new.setScale_in_out(virtualDeploymentUnit.getScale_in_out());
  HashSet<String> monitoringParameters = new HashSet<>();
  monitoringParameters.addAll(virtualDeploymentUnit.getMonitoring_parameter());
  vdu_new.setMonitoring_parameter(monitoringParameters);
  vdu_new.setVdu_constraint(virtualDeploymentUnit.getVdu_constraint());

  //Set Faultmanagement policies
  Set<VRFaultManagementPolicy> vrFaultManagementPolicies = new HashSet<>();
  if (virtualDeploymentUnit.getFault_management_policy() != null) {
    log.debug(
        "Adding the fault management policies: "
            + virtualDeploymentUnit.getFault_management_policy());
    for (VRFaultManagementPolicy vrfmp : virtualDeploymentUnit.getFault_management_policy()) {
      vrFaultManagementPolicies.add(vrfmp);
    }
  }
  vdu_new.setFault_management_policy(vrFaultManagementPolicies);
  //Set Faultmanagement policies end

  HashSet<String> vmImages = new HashSet<>();
  vmImages.addAll(virtualDeploymentUnit.getVm_image());
  vdu_new.setVm_image(vmImages);

  vdu_new.setVirtual_network_bandwidth_resource(
      virtualDeploymentUnit.getVirtual_network_bandwidth_resource());
  vdu_new.setVirtual_memory_resource_element(
      virtualDeploymentUnit.getVirtual_memory_resource_element());
  virtualNetworkFunctionRecord.getVdu().add(vdu_new);


  return &newVDU
}
