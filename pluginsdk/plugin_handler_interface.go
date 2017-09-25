/*
 *  Copyright (c) 2017 Open Baton (http://openbaton.org)
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package pluginsdk

import (
	"github.com/openbaton/go-openbaton/catalogue"
)



// The Handler interface defines an abstraction of the operations that a VNFM should provide.
type HandlerVim interface {
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
		network []*catalogue.VNFDConnectionPoint,
		secGroup []string,
		userData string) (*catalogue.Server, error)

	LaunchInstanceAndWait(
		vimInstance *catalogue.VIMInstance,
		hostname, image, extID, keyPair string,
		network []*catalogue.VNFDConnectionPoint,
		securityGroups []string,
		s string) (*catalogue.Server, error)

	LaunchInstanceAndWaitWithIPs(
		vimInstance *catalogue.VIMInstance,
		hostname, image, extID, keyPair string,
		network []*catalogue.VNFDConnectionPoint,
		securityGroups []string,
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
