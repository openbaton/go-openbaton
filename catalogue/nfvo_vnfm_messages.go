package catalogue

type VNFMAllocateResourcesMessage struct {
	VNFR         *VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VIMInstances map[string]*VIMInstance       `json:"vimInstances"`
	Userdata     string                        `json:"userdata"`
	keyPairs     []*Key                        `json:"keyPairs"`
}

func (*VNFMAllocateResourcesMessage) Action() Action {
	return ActionAllocateResources
}

type VNFMErrorMessage struct {
	NSRID     string                        `json:"nsrId"`
	VNFR      *VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	Exception map[string]interface{}        `json:"exception"` // I don't know how to deserialize a Java exception
}

func (*VNFMErrorMessage) Action() Action {
	return ActionError
}
