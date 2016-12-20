package messages

import "github.com/mcilloni/go-openbaton/catalogue"

type VNFMAllocateResources struct {
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VIMInstances map[string]*catalogue.VIMInstance       `json:"vimInstances"`
	Userdata     string                                  `json:"userdata"`
	KeyPairs     []*catalogue.Key                        `json:"keyPairs"`
}

type VNFMError struct {
	NSRID     string                                  `json:"nsrId"`
	VNFR      *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	Exception map[string]interface{}                  `json:"exception"` // I don't know how to deserialize a Java exception
}

type VNFMGeneric struct {
	VNFR                *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VNFRecordDependency *catalogue.VNFRecordDependency          `json:"vnfRecordDependency"`
}

type VNFMGrantLifecycleOperation struct {
	VNFD                 *catalogue.VirtualNetworkFunctionDescriptor `json:"virtualNetworkFunctionDescriptor"`
	VDUSet               []*catalogue.VirtualDeploymentUnit          `json:"vduSet"`
	DeploymentFlavourKey string                                      `json:"deploymentFlavourKey"`
}

type VNFMHealed struct {
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance"`
	Cause        string                                  `json:"cause"`
}

type VNFMInstantiate struct {
	VNFR *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
}

type VNFMScaled struct {
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	vnfcInstance *catalogue.VNFCInstance                 `json:"vnfcInstance"`
}

type VNFMScaling struct {
	VNFR     *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	userData string                                  `json:"userData"`
}

type VNFMStartStop struct {
	VNFR           *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VNFCInstance   *catalogue.VNFCInstance                 `json:"vnfcInstance"`
	VNFRDependency *catalogue.VNFRecordDependency          `json:"vnfrDependency"`
}
