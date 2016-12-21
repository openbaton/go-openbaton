package messages

import (
	"github.com/mcilloni/go-openbaton/catalogue"
)

type OrAllocateResources struct {
	VDUSet []*catalogue.VirtualDeploymentUnit `json:"vduSet"`
}

type OrError struct {
	VNFR    *catalogue.VirtualNetworkFunctionRecord `json:"vnfr"`
	Message string                                  `json:"message"`
}

type OrGeneric struct {
	VNFR           *catalogue.VirtualNetworkFunctionRecord `json:"vnfr"`
	VNFRDependency *catalogue.VNFRecordDependency          `json:"vnfrd"`
}

type OrGrantLifecycleOperation struct {
	GrantAllowed bool                                    `json:"grantAllowed"`
	VDUVIM       map[string]*catalogue.VIMInstance       `json:"vduVim"`
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
}

type OrHealVNFRequest struct {
	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance"`
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	Cause        string                                  `json:"cause"`
}

type OrInstantiate struct {
	VNFD            *catalogue.VirtualNetworkFunctionDescriptor `json:"vnfd"`
	VNFDFlavour     *catalogue.VNFDeploymentFlavour             `json:"vnfdf"`
	VNFInstanceName string                                      `json:"vnfInstanceName"`
	VLRs            []*catalogue.VirtualLinkRecord              `json:"vlrs"`
	Extension       map[string]string                           `json:"extension"`
	VIMInstances    map[string][]*catalogue.VIMInstance         `json:"vimInstances"`
	VNFPackage      *catalogue.VNFPackage                       `json:"vnfPackage"`
	Keys            []*catalogue.Key                            `json:"keys"`
}

type OrScaling struct {
	Component    *catalogue.VNFComponent                 `json:"component"`
	VNFCInstance *catalogue.VNFCInstance                 `json:"vnfcInstance"`
	VIMInstance  *catalogue.VIMInstance                  `json:"vimInstance"`
	VNFPackage   *catalogue.VNFPackage                   `json:"vnfPackage"`
	VNFR         *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	Dependency   *catalogue.VNFRecordDependency          `json:"dependency"`
	Mode         string                                  `json:"mode"`
	Extension    map[string]string                       `json:"extension"`
}

type OrStartStop struct {
	VNFR           *catalogue.VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VNFCInstance   *catalogue.VNFCInstance                 `json:"vnfcInstance"`
	VNFRDependency *catalogue.VNFRecordDependency          `json:"vnfrDependency"`
}

type OrUpdate struct {
	Script *catalogue.Script                       `json:"script"`
	VNFR   *catalogue.VirtualNetworkFunctionRecord `json:"vnfr"`
}
