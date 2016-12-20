package catalogue

type vnfmAllocateResourcesMessage struct {
	VNFR         *VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VIMInstances map[string]*VIMInstance       `json:"vimInstances"`
	Userdata     string                        `json:"userdata"`
	KeyPairs     []*Key                        `json:"keyPairs"`
}

type vnfmErrorMessage struct {
	NSRID     string                        `json:"nsrId"`
	VNFR      *VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	Exception map[string]interface{}        `json:"exception"` // I don't know how to deserialize a Java exception
}

type vnfmGenericMessage struct {
	VNFR                *VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VNFRecordDependency *VNFRecordDependency          `json:"vnfRecordDependency"`
}

type vnfmGrantLifecycleOperationMessage struct {
	virtualNetworkFunctionDescriptor *VirtualNetworkFunctionDescriptor
	vduSet                           []*VirtualDeploymentUnit
	deploymentFlavourKey             string
}

type vnfmOrHealedMessage struct {
	virtualNetworkFunctionRecord *VirtualNetworkFunctionRecord
	vnfcInstance                 *VNFCInstance
	cause                        string
}
