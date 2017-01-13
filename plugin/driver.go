package plugin

import (
	"github.com/mcilloni/go-openbaton/catalogue"
)

const interfaceVersion = "1.0"

// Driver describes a VIMDriver.
// Every driver implementation must adhere to this interface and implements its methods.
type Driver interface {
	AddFlavour(vimInstance *catalogue.VIMInstance, deploymentFlavour *catalogue.DeploymentFlavour) (*catalogue.DeploymentFlavour, error)

	AddImage(vimInstance *catalogue.VIMInstance, image *catalogue.NFVImage, imageFile []byte) (*catalogue.NFVImage, error)

	AddImageFromURL(vimInstance *catalogue.VIMInstance, image *catalogue.NFVImage, imageURL string) (*catalogue.NFVImage, error)

	CopyImage(vimInstance *catalogue.VIMInstance, image *catalogue.NFVImage, imageFile []byte) (*catalogue.NFVImage, error)

	CreateNetwork(vimInstance *catalogue.VIMInstance, network *catalogue.Network) (*catalogue.Network, error)

	CreateSubnet(vimInstance *catalogue.VIMInstance, createdNetwork *catalogue.Network, subnet *catalogue.Subnet) (*catalogue.Subnet, error)

	DeleteFlavour(vimInstance *catalogue.VIMInstance, extID string) (bool, error)

	DeleteImage(vimInstance *catalogue.VIMInstance, image *catalogue.NFVImage) (bool, error)

	DeleteNetwork(vimInstance *catalogue.VIMInstance, extID string) (bool, error)

	DeleteServerByIDAndWait(vimInstance *catalogue.VIMInstance, id string) error

	DeleteSubnet(vimInstance *catalogue.VIMInstance, existingSubnetExtID string) (bool, error)

	LaunchInstance(
		vimInstance *catalogue.VIMInstance,
		name, image, Flavour, keypair string,
		network, secGroup []string,
		userData string) (*catalogue.Server, error)

	LaunchInstanceAndWait(
		vimInstance *catalogue.VIMInstance,
		hostname, image, extID, keyPair string,
		networks, securityGroups []string,
		s string) (*catalogue.Server, error)

	LaunchInstanceAndWaitWithIPs(
		vimInstance *catalogue.VIMInstance,
		hostname, image, extID, keyPair string,
		networks, securityGroups []string,
		s string,
		floatingIps map[string]string,
		keys []*catalogue.Key) (*catalogue.Server, error)

	ListFlavours(vimInstance *catalogue.VIMInstance) ([]*catalogue.DeploymentFlavour, error)

	ListImages(vimInstance *catalogue.VIMInstance) ([]*catalogue.NFVImage, error)

	ListNetworks(vimInstance *catalogue.VIMInstance) ([]*catalogue.Network, error)

	ListServer(vimInstance *catalogue.VIMInstance) ([]*catalogue.Server, error)

	NetworkByID(vimInstance *catalogue.VIMInstance, id string) (*catalogue.Network, error)

	Quota(vimInstance *catalogue.VIMInstance) (*catalogue.Quota, error)

	SubnetsExtIDs(vimInstance *catalogue.VIMInstance, networkExtID string) ([]string, error)

	Type(vimInstance *catalogue.VIMInstance) (string, error)

	UpdateFlavour(vimInstance *catalogue.VIMInstance, deploymentFlavour *catalogue.DeploymentFlavour) (*catalogue.DeploymentFlavour, error)

	UpdateImage(vimInstance *catalogue.VIMInstance, image *catalogue.NFVImage) (*catalogue.NFVImage, error)

	UpdateNetwork(vimInstance *catalogue.VIMInstance, network *catalogue.Network) (*catalogue.Network, error)

	UpdateSubnet(vimInstance *catalogue.VIMInstance, createdNetwork *catalogue.Network, subnet *catalogue.Subnet) (*catalogue.Subnet, error)
}

// DriverError is a special error type that also specifies a catalogue.Server
// to be returned to the NFVO.
type DriverError struct {
	Message           string `json:"detailMessage"`
	*catalogue.Server `json:"server"`
}

// Error returns a description of the error.
func (e DriverError) Error() string {
	return e.Message + " on server " + e.Server.Name
}
