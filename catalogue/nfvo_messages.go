package catalogue

type NFVMessage interface {
	Action() Action
}

type VNFMAllocateResourcesMessage struct {
	VNFR         *VirtualNetworkFunctionRecord `json:"virtualNetworkFunctionRecord"`
	VIMInstances map[string]*VIMInstance       `json:"vimInstances"`
	Userdata     string                        `json:"userdata"`
	keyPairs     []*Key                        `json:"keyPairs"`
}

func (*VNFMAllocateResourcesMessage) Action() Action {
	return ActionAllocateResources
}
