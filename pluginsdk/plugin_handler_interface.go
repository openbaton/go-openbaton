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
	AddFlavour(vimInstance interface{}, deploymentFlavour *catalogue.DeploymentFlavour) (*catalogue.DeploymentFlavour, error)

	AddImage(vimInstance interface{}, image *catalogue.DockerImage, imageFile []byte) (*catalogue.DockerImage, error)

	AddImageFromURL(vimInstance interface{}, image *catalogue.DockerImage, imageURL string) (*catalogue.DockerImage, error)

	CopyImage(vimInstance interface{}, image *catalogue.DockerImage, imageFile []byte) (*catalogue.DockerImage, error)

	CreateNetwork(vimInstance interface{}, network *catalogue.DockerNetwork) (*catalogue.DockerNetwork, error)

	CreateSubnet(vimInstance interface{}, createdNetwork *catalogue.DockerNetwork, subnet *catalogue.Subnet) (*catalogue.Subnet, error)

	DeleteFlavour(vimInstance interface{}, extID string) (bool, error)

	DeleteImage(vimInstance interface{}, image *catalogue.DockerImage) (bool, error)

	DeleteNetwork(vimInstance interface{}, extID string) (bool, error)

	DeleteServerByIDAndWait(vimInstance interface{}, id string) error

	DeleteSubnet(vimInstance interface{}, existingSubnetExtID string) (bool, error)

	Refresh(vimInstance interface{}) (interface{}, error)

	LaunchInstance(
		vimInstance interface{},
		name, image, Flavour, keypair string,
		network []*catalogue.VNFDConnectionPoint,
		secGroup []string,
		userData string) (*catalogue.Server, error)

	LaunchInstanceAndWait(
		vimInstance interface{},
		hostname, image, extID, keyPair string,
		network []*catalogue.VNFDConnectionPoint,
		securityGroups []string,
		s string) (*catalogue.Server, error)

	LaunchInstanceAndWaitWithIPs(
		vimInstance interface{},
		hostname, image, extID, keyPair string,
		network []*catalogue.VNFDConnectionPoint,
		securityGroups []string,
		s string,
		floatingIps map[string]string,
		keys []*catalogue.Key) (*catalogue.Server, error)

	ListFlavours(vimInstance interface{}) ([]*catalogue.DeploymentFlavour, error)

	ListImages(vimInstance interface{}) ([]*catalogue.DockerImage, error)

	ListNetworks(vimInstance interface{}) ([]*catalogue.DockerNetwork, error)

	ListServer(vimInstance interface{}) ([]*catalogue.Server, error)

	NetworkByID(vimInstance interface{}, id string) (*catalogue.DockerNetwork, error)

	Quota(vimInstance interface{}) (*catalogue.Quota, error)

	SubnetsExtIDs(vimInstance interface{}, networkExtID string) ([]string, error)

	Type(vimInstance interface{}) (string, error)

	UpdateFlavour(vimInstance interface{}, deploymentFlavour *catalogue.DeploymentFlavour) (*catalogue.DeploymentFlavour, error)

	UpdateImage(vimInstance interface{}, image *catalogue.DockerImage) (*catalogue.DockerImage, error)

	UpdateNetwork(vimInstance interface{}, network *catalogue.DockerNetwork) (*catalogue.DockerNetwork, error)

	UpdateSubnet(vimInstance interface{}, createdNetwork *catalogue.DockerNetwork, subnet *catalogue.Subnet) (*catalogue.Subnet, error)
}
