package plugin

import (
	"github.com/mcilloni/go-openbaton/catalogue"
)

const interfaceVersion = "1.0"

type Driver interface {
	AddFlavor(vimInstance *catalogue.VIMInstance, deploymentFlavour *catalogue.DeploymentFlavour) (*catalogue.DeploymentFlavour, error)

	// imageData can be both a string containing an URL OR a byte slice containing an image
	AddImage(vimInstance *catalogue.VIMInstance, image *catalogue.NFVImage, imageData interface{}) (*catalogue.NFVImage, error)

	CopyImage(vimInstance *catalogue.VIMInstance, image *catalogue.NFVImage, imageFile []byte) (*catalogue.NFVImage, error)

	CreateNetwork(vimInstance *catalogue.VIMInstance, network *catalogue.Network) (*catalogue.Network, error)

	CreateSubnet(vimInstance *catalogue.VIMInstance, createdNetwork *catalogue.Network, subnet *catalogue.Subnet) (*catalogue.Subnet, error)

	DeleteFlavor(vimInstance *catalogue.VIMInstance, extID string) (bool, error)

	DeleteImage(vimInstance *catalogue.VIMInstance, image *catalogue.NFVImage) (bool, error)

	DeleteNetwork(vimInstance *catalogue.VIMInstance, extID string) (bool, error)

	DeleteServerByIDAndWait(vimInstance *catalogue.VIMInstance, id string) error

	DeleteSubnet(vimInstance *catalogue.VIMInstance, existingSubnetExtID string) (bool, error)

	GetNetworkById(vimInstance *catalogue.VIMInstance, id string) (*catalogue.Network, error)

	GetQuota(vimInstance *catalogue.VIMInstance) (*catalogue.Quota, error)

	GetSubnetsExtIDs(vimInstance *catalogue.VIMInstance, networkExtID string) ([]string, error)

	GetType(vimInstance *catalogue.VIMInstance) (string, error)

	LaunchInstance(
		vimInstance *catalogue.VIMInstance,
		name, image, flavor, keypair string,
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

	ListFlavors(vimInstance *catalogue.VIMInstance) ([]*catalogue.DeploymentFlavour, error)

	ListImages(vimInstance *catalogue.VIMInstance) ([]*catalogue.NFVImage, error)

	ListNetworks(vimInstance *catalogue.VIMInstance) ([]*catalogue.Network, error)

	ListServer(vimInstance *catalogue.VIMInstance) ([]*catalogue.Server, error)

	UpdateFlavor(vimInstance *catalogue.VIMInstance, deploymentFlavour *catalogue.DeploymentFlavour) (*catalogue.DeploymentFlavour, error)

	UpdateImage(vimInstance *catalogue.VIMInstance, image *catalogue.NFVImage) (*catalogue.NFVImage, error)

	UpdateNetwork(vimInstance *catalogue.VIMInstance, network *catalogue.Network) (*catalogue.Network, error)

	UpdateSubnet(vimInstance *catalogue.VIMInstance, createdNetwork *catalogue.Network, subnet *catalogue.Subnet) (*catalogue.Subnet, error)
}
