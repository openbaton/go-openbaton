package messages

import "github.com/mcilloni/go-openbaton/catalogue"

type vnfmAllocateResourcesMessage struct {
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VIMInstances map[string]*catalogue.VIMInstance       `json:"vimInstances"`
	Userdata     string                        `json:"userdata"`
	KeyPairs     []*catalogue.Key                        `json:"keyPairs"`
}

type vnfmErrorMessage struct {
	NSRID     string                        `json:"nsrId"`
	VNFR      *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	Exception map[string]interface{}        `json:"exception"` // I don't know how to deserialize a Java exception
}

type vnfmGenericMessage struct {
	VNFR                *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VNFRecordDependency *catalogue.VNFRecordDependency          `json:"vnfRecordDependency"`
}

type vnfmGrantLifecycleOperationMessage struct {
	virtualNetworkFunctionDescriptor *catalogue.VirtualNetworkFunctionDescriptor
	vduSet                           []*catalogue.VirtualDeploymentUnit
	deploymentFlavourKey             string
}

type vnfmHealedMessage struct {
	virtualNetworkFunctionRecord *catalogue.VirtualNetworkFunctionRecord
	vnfcInstance                 *catalogue.VNFCInstance
	cause                        string
}

type vnfmInstantiateMessage struct {
	
}