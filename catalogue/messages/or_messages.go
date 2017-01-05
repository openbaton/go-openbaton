package messages

import (
	"github.com/mcilloni/go-openbaton/catalogue"
)

type orMessage struct{}

func (orMessage) From() SenderType {
	return NFVO
}

type OrAllocateResources struct {
	orMessage

	VDUSet []*catalogue.VirtualDeploymentUnit `json:"vduSet,omitempty"`
}

func (OrAllocateResources) DefaultAction() catalogue.Action {
	return catalogue.ActionAllocateResources
}

type OrError struct {
	orMessage

	VNFR    *catalogue.VirtualNetworkFunctionRecord `json:"vnfr,omitempty"`
	Message string                                  `json:"message,omitempty"`
}

func (OrError) DefaultAction() catalogue.Action {
	return catalogue.ActionError
}

type OrGeneric struct {
	orMessage

	VNFR           *catalogue.VirtualNetworkFunctionRecord `json:"vnfr,omitempty"`
	VNFRDependency *catalogue.VNFRecordDependency          `json:"vnfrd,omitempty"`
}

func (OrGeneric) DefaultAction() catalogue.Action {
	return catalogue.NoActionSpecified
}

type OrGrantLifecycleOperation struct {
	orMessage

	GrantAllowed bool                                    `json:"grantAllowed"`
	VDUVIM       map[string]*catalogue.VIMInstance       `json:"vduVim,omitempty"`
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
}

func (OrGrantLifecycleOperation) DefaultAction() catalogue.Action {
	return catalogue.ActionGrantOperation
}

type OrHealVNFRequest struct {
	orMessage

	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance,omitempty"`
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	Cause        string                                  `json:"cause,omitempty"`
}

func (OrHealVNFRequest) DefaultAction() catalogue.Action {
	return catalogue.ActionHeal
}

type OrInstantiate struct {
	orMessage

	VNFD            *catalogue.VirtualNetworkFunctionDescriptor `json:"vnfd,omitempty"`
	VNFDFlavour     *catalogue.VNFDeploymentFlavour             `json:"vnfdf,omitempty"`
	VNFInstanceName string                                      `json:"vnfInstanceName,omitempty"`
	VLRs            []*catalogue.VirtualLinkRecord              `json:"vlrs,omitempty"`
	Extension       map[string]string                           `json:"extension,omitempty"`
	VIMInstances    map[string][]*catalogue.VIMInstance   `json:"vimInstances,omitempty"`
	VNFPackage      *catalogue.VNFPackage                       `json:"vnfPackage,omitempty"`
	Keys            []*catalogue.Key                            `json:"keys,omitempty"`
}

func (OrInstantiate) DefaultAction() catalogue.Action {
	return catalogue.ActionInstantiate
}

type OrScaling struct {
	orMessage

	Component    *catalogue.VNFComponent                 `json:"component,omitempty"`
	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance,omitempty"`
	VIMInstance  *catalogue.VIMInstance                  `json:"vimInstance,omitempty"`
	VNFPackage   *catalogue.VNFPackage                   `json:"vnfPackage,omitempty"`
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	Dependency   *catalogue.VNFRecordDependency          `json:"dependency,omitempty"`
	Mode         string                                  `json:"mode,omitempty"`
	Extension    map[string]string                       `json:"extension,omitempty"`
}

func (OrScaling) DefaultAction() catalogue.Action {
	return catalogue.ActionScaling
}

type OrStartStop struct {
	orMessage

	VNFR           *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	VNFCInstance   *catalogue.VNFCInstance                 `json:"vnfcInstance,omitempty"`
	VNFRDependency *catalogue.VNFRecordDependency          `json:"vnfrDependency,omitempty"`
}

func (OrStartStop) DefaultAction() catalogue.Action {
	return catalogue.ActionStart
}

type OrUpdate struct {
	orMessage

	Script *catalogue.Script                       `json:"script,omitempty"`
	VNFR   *catalogue.VirtualNetworkFunctionRecord `json:"vnfr,omitempty"`
}

func (OrUpdate) DefaultAction() catalogue.Action {
	return catalogue.ActionUpdate
}

