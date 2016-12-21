package messages

import "github.com/mcilloni/go-openbaton/catalogue"

type vnfmMessage struct{}

func (vnfmMessage) From() SenderType {
	return VNFM
}

type VNFMAllocateResources struct {
	vnfmMessage

	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VIMInstances map[string]*catalogue.VIMInstance       `json:"vimInstances"`
	Userdata     string                                  `json:"userdata"`
	KeyPairs     []*catalogue.Key                        `json:"keyPairs"`
}

func (VNFMAllocateResources) DefaultAction() catalogue.Action {
	return catalogue.ActionAllocateResources
}

type VNFMError struct {
	vnfmMessage

	NSRID     string                                  `json:"nsrId"`
	VNFR      *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	Exception map[string]interface{}                  `json:"exception"` // I don't know how to deserialize a Java exception
}

func (VNFMError) DefaultAction() catalogue.Action {
	return catalogue.ActionError
}

type VNFMGeneric struct {
	vnfmMessage

	VNFR                *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VNFRecordDependency *catalogue.VNFRecordDependency          `json:"vnfRecordDependency"`
}

func (VNFMGeneric) DefaultAction() catalogue.Action {
	return catalogue.NoActionSpecified
}

type VNFMGrantLifecycleOperation struct {
	vnfmMessage

	VNFD                 *catalogue.VirtualNetworkFunctionDescriptor `json:"virtualNetworkFunctionDescriptor"`
	VDUSet               []*catalogue.VirtualDeploymentUnit          `json:"vduSet"`
	DeploymentFlavourKey string                                      `json:"deploymentFlavourKey"`
}

func (VNFMGrantLifecycleOperation) DefaultAction() catalogue.Action {
	return catalogue.ActionGrantOperation
}

type VNFMHealed struct {
	vnfmMessage

	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance"`
	Cause        string                                  `json:"cause"`
}

func (VNFMHealed) DefaultAction() catalogue.Action {
	return catalogue.ActionHeal
}

type VNFMInstantiate struct {
	vnfmMessage

	VNFR *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
}

func (VNFMInstantiate) DefaultAction() catalogue.Action {
	return catalogue.ActionInstantiate
}

type VNFMScaled struct {
	vnfmMessage

	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	vnfcInstance *catalogue.VNFCInstance                 `json:"vnfcInstance"`
}

func (VNFMScaled) DefaultAction() catalogue.Action {
	return catalogue.ActionScaled
}

type VNFMScaling struct {
	vnfmMessage

	VNFR     *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	userData string                                  `json:"userData"`
}

func (VNFMScaling) DefaultAction() catalogue.Action {
	return catalogue.ActionScaling
}

type VNFMStartStop struct {
	vnfmMessage

	VNFR           *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VNFCInstance   *catalogue.VNFCInstance                 `json:"vnfcInstance"`
	VNFRDependency *catalogue.VNFRecordDependency          `json:"vnfrDependency"`
}

func (VNFMStartStop) DefaultAction() catalogue.Action {
	return catalogue.ActionStart
}
