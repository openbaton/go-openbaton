package messages

import (
	"github.com/mcilloni/go-openbaton/catalogue"
)

//go:generate stringer -type=orMessage
type orMessage struct{}

func (orMessage) From() SenderType {
	return NFVO
}

//go:generate stringer -type=OrAllocateResources
type OrAllocateResources struct {
	orMessage

	VDUSet []*catalogue.VirtualDeploymentUnit `json:"vduSet,omitempty"`
}

func (OrAllocateResources) DefaultAction() catalogue.Action {
	return catalogue.ActionAllocateResources
}

//go:generate stringer -type=OrError
type OrError struct {
	orMessage

	VNFR    *catalogue.VirtualNetworkFunctionRecord `json:"vnfr,omitempty"`
	Message string                                  `json:"message,omitempty"`
}

func (OrError) DefaultAction() catalogue.Action {
	return catalogue.ActionError
}

//go:generate stringer -type=OrGeneric
type OrGeneric struct {
	orMessage

	VNFR           *catalogue.VirtualNetworkFunctionRecord `json:"vnfr,omitempty"`
	VNFRDependency *catalogue.VNFRecordDependency          `json:"vnfrd,omitempty"`
}

func (OrGeneric) DefaultAction() catalogue.Action {
	return catalogue.NoActionSpecified
}

//go:generate stringer -type=OrGrantLifecycleOperation
type OrGrantLifecycleOperation struct {
	orMessage

	GrantAllowed bool                                    `json:"grantAllowed"`
	VDUVIM       map[string]*catalogue.VIMInstance       `json:"vduVim,omitempty"`
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
}

func (OrGrantLifecycleOperation) DefaultAction() catalogue.Action {
	return catalogue.ActionGrantOperation
}

//go:generate stringer -type=OrHealVNFRequest
type OrHealVNFRequest struct {
	orMessage

	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance,omitempty"`
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	Cause        string                                  `json:"cause,omitempty"`
}

func (OrHealVNFRequest) DefaultAction() catalogue.Action {
	return catalogue.ActionHeal
}

//go:generate stringer -type=OrInstantiate
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

//go:generate stringer -type=OrScaling
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

//go:generate stringer -type=OrStartStop
type OrStartStop struct {
	orMessage

	VNFR           *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord,omitempty"`
	VNFCInstance   *catalogue.VNFCInstance                 `json:"vnfcInstance,omitempty"`
	VNFRDependency *catalogue.VNFRecordDependency          `json:"vnfrDependency,omitempty"`
}

func (OrStartStop) DefaultAction() catalogue.Action {
	return catalogue.ActionStart
}

//go:generate stringer -type=OrUpdate
type OrUpdate struct {
	orMessage

	Script *catalogue.Script                       `json:"script,omitempty"`
	VNFR   *catalogue.VirtualNetworkFunctionRecord `json:"vnfr,omitempty"`
}

func (OrUpdate) DefaultAction() catalogue.Action {
	return catalogue.ActionUpdate
}
