package openbaton

// An extended Virtual Link based on ETSI GS NFV-MAN 001 V1.1.1 (2014-12)
type InternalVirtualLink struct {
	VirtualLink
	ConnectionPointsReferences []string `json:"connection_points_references"`
}

type NetworkForwardingPath struct {
	ID string `json:"id"`
	Version int `json:"version"` 
	Policy *Policy `json:"policy"`
	Connection map[string]string `json:"connection"`
}