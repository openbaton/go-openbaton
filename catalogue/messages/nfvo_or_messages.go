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

	VDUSet []*catalogue.VirtualDeploymentUnit `json:"vduSet"`
}

func (OrAllocateResources) DefaultAction() catalogue.Action {
	return catalogue.ActionAllocateResources
}

type OrError struct {
	orMessage

	VNFR    *catalogue.VirtualNetworkFunctionRecord `json:"vnfr"`
	Message string                                  `json:"message"`
}

func (OrError) DefaultAction() catalogue.Action {
	return catalogue.ActionError
}

type OrGeneric struct {
	orMessage

	VNFR           *catalogue.VirtualNetworkFunctionRecord `json:"vnfr"`
	VNFRDependency *catalogue.VNFRecordDependency          `json:"vnfrd"`
}

func (OrGeneric) DefaultAction() catalogue.Action {
	return catalogue.NoActionSpecified
}

type OrGrantLifecycleOperation struct {
	orMessage

	GrantAllowed bool                                    `json:"grantAllowed"`
	VDUVIM       map[string]*catalogue.VIMInstance       `json:"vduVim"`
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
}

func (OrGrantLifecycleOperation) DefaultAction() catalogue.Action {
	return catalogue.ActionGrantOperation
}

type OrHealVNFRequest struct {
	orMessage

	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance"`
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	Cause        string                                  `json:"cause"`
}

func (OrHealVNFRequest) DefaultAction() catalogue.Action {
	return catalogue.ActionHeal
}

type OrInstantiate struct {
	orMessage

	VNFD            *catalogue.VirtualNetworkFunctionDescriptor `json:"vnfd"`
	VNFDFlavour     *catalogue.VNFDeploymentFlavour             `json:"vnfdf"`
	VNFInstanceName string                                      `json:"vnfInstanceName"`
	VLRs            []*catalogue.VirtualLinkRecord              `json:"vlrs"`
	Extension       map[string]string                           `json:"extension"`
	VIMInstances    map[string][]*catalogue.VIMInstance         `json:"vimInstances"`
	VNFPackage      *catalogue.VNFPackage                       `json:"vnfPackage"`
	Keys            []*catalogue.Key                            `json:"keys"`
}

func (OrInstantiate) DefaultAction() catalogue.Action {
	return catalogue.ActionInstantiate
}

type OrScaling struct {
	orMessage

	Component    *catalogue.VNFComponent                 `json:"component"`
	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance"`
	VIMInstance  *catalogue.VIMInstance                  `json:"vimInstance"`
	VNFPackage   *catalogue.VNFPackage                   `json:"vnfPackage"`
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	Dependency   *catalogue.VNFRecordDependency          `json:"dependency"`
	Mode         string                                  `json:"mode"`
	Extension    map[string]string                       `json:"extension"`
}

func (OrScaling) DefaultAction() catalogue.Action {
	return catalogue.ActionScaling
}

type OrStartStop struct {
	orMessage

	VNFR           *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VNFCInstance   *catalogue.VNFCInstance                 `json:"vnfcInstance"`
	VNFRDependency *catalogue.VNFRecordDependency          `json:"vnfrDependency"`
}

func (OrStartStop) DefaultAction() catalogue.Action {
	return catalogue.ActionStart
}

type OrUpdate struct {
	orMessage

	Script *catalogue.Script                       `json:"script"`
	VNFR   *catalogue.VirtualNetworkFunctionRecord `json:"vnfr"`
}

func (OrUpdate) DefaultAction() catalogue.Action {
	return catalogue.ActionUpdate
}
