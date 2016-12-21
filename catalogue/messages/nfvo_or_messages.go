package messages

import (
	"github.com/mcilloni/go-openbaton/catalogue"
)

type OrAllocateResources struct {
    VDUSet []*catalogue.VirtualDeploymentUnit `json:"vduSet"`
}

type OrError struct {
    vnfr *catalogue.VirtualNetworkFunctionRecord
    message string
}

type OrGeneric struct {
    vnfr *catalogue.VirtualNetworkFunctionRecord
    vnfrd *catalogue.VNFRecordDependency
}

type OrGrantLifecycleOperation struct {
    grantAllowed bool
    vduVim map[string]*catalogue.VIMInstance
    virtualNetworkFunctionRecord *catalogue.VirtualNetworkFunctionRecord
}

type OrHealVNFRequest struct {
    vnfcInstance *catalogue.VNFCInstance
    virtualNetworkFunctionRecord *catalogue.VirtualNetworkFunctionRecord
    cause string
}

type OrInstantiate struct {
    vnfd *catalogue.VirtualNetworkFunctionDescriptor
    vnfdf *catalogue.VNFDeploymentFlavour
    vnfInstanceName string
    vlrs []*catalogue.VirtualLinkRecord
    extension map[string]string
    vimInstances map[string][]*catalogue.VIMInstance
    vnfPackage *catalogue.VNFPackage
    keys []*catalogue.Key
}

type OrScaling struct {
    component *catalogue.VNFComponent
    vnfcInstance *catalogue.VNFCInstance
    vimInstance *catalogue.VimInstance
    vnfPackage *catalogue.VNFPackage
    virtualNetworkFunctionRecord *catalogue.VirtualNetworkFunctionRecord
    dependency *catalogue.VNFRecordDependency
    mode string
    extension map[string]string
}

type OrStartStop struct {
    virtualNetworkFunctionRecord *catalogue.VirtualNetworkFunctionRecord
    vnfcInstance *catalogue.VNFCInstance
    vnfrd *catalogue.VNFRecordDependency
    vnfrDependency *catalogue.VNFRecordDependency
}

type OrUpdate struct {
    script *catalogue.Script
    vnfr *catalogue.VirtualNetworkFunctionRecord
}
