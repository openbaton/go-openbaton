package messages

import "github.com/openbaton/go-openbaton/catalogue"

type vnfmMessage struct{}

func (vnfmMessage) From() SenderType {
	return VNFM
}

type VNFMAllocateResources struct {
	vnfmMessage

	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	VIMInstances map[string]*catalogue.VIMInstance       `json:"vimInstances,omitempty"`
	Userdata     string                                  `json:"userdata,omitempty"`
	KeyPairs     []*catalogue.Key                        `json:"keyPairs,omitempty"`
}

func (VNFMAllocateResources) DefaultAction() catalogue.Action {
	return catalogue.ActionAllocateResources
}

type VNFMError struct {
	vnfmMessage

	NSRID     string                                  `json:"nsrId,omitempty"`
	VNFR      *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	Exception map[string]interface{}                  `json:"exception,omitempty"` // I don't know how to deserialize a Java exception
}

func (VNFMError) DefaultAction() catalogue.Action {
	return catalogue.ActionError
}

type VNFMGeneric struct {
	vnfmMessage

	VNFR                *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	VNFRecordDependency *catalogue.VNFRecordDependency          `json:"vnfRecordDependency,omitempty"`
}

func (VNFMGeneric) DefaultAction() catalogue.Action {
	return catalogue.NoActionSpecified
}

type VNFMGrantLifecycleOperation struct {
	vnfmMessage

	VNFD                 *catalogue.VirtualNetworkFunctionDescriptor `json:"virtualNetworkFunctionDescriptor,omitempty"`
	VDUSet               []*catalogue.VirtualDeploymentUnit          `json:"vduSet,omitempty"`
	DeploymentFlavourKey string                                      `json:"deploymentFlavourKey,omitempty"`
}

func (VNFMGrantLifecycleOperation) DefaultAction() catalogue.Action {
	return catalogue.ActionGrantOperation
}

type VNFMHealed struct {
	vnfmMessage

	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance,omitempty"`
	Cause        string                                  `json:"cause,omitempty"`
}

func (VNFMHealed) DefaultAction() catalogue.Action {
	return catalogue.ActionHeal
}

type VNFMInstantiate struct {
	vnfmMessage

	VNFR *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
}

func (VNFMInstantiate) DefaultAction() catalogue.Action {
	return catalogue.ActionInstantiate
}

type VNFMScaled struct {
	vnfmMessage

	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance,omitempty"`
}

func (VNFMScaled) DefaultAction() catalogue.Action {
	return catalogue.ActionScaled
}

type VNFMScaling struct {
	vnfmMessage

	VNFR     *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	UserData string                                  `json:"userData,omitempty"`
}

func (VNFMScaling) DefaultAction() catalogue.Action {
	return catalogue.ActionScaling
}

type VNFMStartStop struct {
	vnfmMessage

	VNFR           *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	VNFCInstance   *catalogue.VNFCInstance                 `json:"vnfcInstance,omitempty"`
	VNFRDependency *catalogue.VNFRecordDependency          `json:"vnfrDependency,omitempty"`
}

func (VNFMStartStop) DefaultAction() catalogue.Action {
	return catalogue.ActionStart
}
