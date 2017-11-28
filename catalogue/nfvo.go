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

// Package catalogue contains a partial implementation of the OpenBaton catalogue.
package catalogue

import (
	"time"
)

type Action string
type ImageStatus string

const (
	Active = ImageStatus("ACTIVE")
	Error  = ImageStatus("ERROR")
)

const (
	ActionGrantOperation         = Action("GRANT_OPERATION")
	ActionAllocateResources      = Action("ALLOCATE_RESOURCES")
	ActionScaleIn                = Action("SCALE_IN")
	ActionScaleOut               = Action("SCALE_OUT")
	ActionScaling                = Action("SCALING")
	ActionError                  = Action("ERROR")
	ActionReleaseResources       = Action("RELEASE_RESOURCES")
	ActionInstantiate            = Action("INSTANTIATE")
	ActionModify                 = Action("MODIFY")
	ActionHeal                   = Action("HEAL")
	ActionUpdateVNFR             = Action("UPDATEVNFR")
	ActionUpdate                 = Action("UPDATE")
	ActionScaled                 = Action("SCALED")
	ActionReleaseResourcesFinish = Action("RELEASE_RESOURCES_FINISH")
	ActionInstantiateFinish      = Action("INSTANTIATE_FINISH")
	ActionConfigure              = Action("CONFIGURE")
	ActionStart                  = Action("START")
	ActionStop                   = Action("STOP")
	ActionResume                 = Action("RESUME")

	NoActionSpecified = Action("")
)

type ConfigurationParameter struct {
	ID          string            `json:"id,omitempty"`
	HbVersion   int               `json:"hbVersion,omitempty"`
	ProjectID   string            `json:"projectId"`
	Shared      bool              `json:"shared,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Description string            `json:"description,omitempty"`
	ConfKey     string            `json:"confKey"`
	Value       string            `json:"value,omitempty"`
}

type Configuration struct {
	ID                      string                    `json:"id,omitempty"`
	HbVersion               int                       `json:"hbVersion,omitempty"`
	ProjectID               string                    `json:"projectId"`
	Shared                  bool                      `json:"shared,omitempty"`
	Metadata                map[string]string         `json:"metadata,omitempty"`
	ConfigurationParameters []*ConfigurationParameter `json:"configurationParameters"`
	Name                    string                    `json:"name"`
}

func (cfg *Configuration) Append(p *ConfigurationParameter) {
	cfg.ConfigurationParameters = append(cfg.ConfigurationParameters, p)
}

type Date string

func NewDate() Date {
	return NewDateWithTime(time.Now())
}

func NewDateWithTime(t time.Time) Date {
	return Date(t.Format("Jan 2, 2006 3:4:5 PM"))
}

func UnixDate(timestamp int64) Date {
	return NewDateWithTime(time.Unix(timestamp, 0))
}

func (d Date) Time() time.Time {
	if d == "" {
		return time.Unix(0, 0)
	}

	t, err := time.Parse("Jan 2, 2006 3:4:5 PM", string(d))
	// Try other formats
	if err != nil {
		t, err = time.Parse("Jan 02, 2006 03:04:05 PM", string(d))
		if err != nil {
			return time.Unix(0, 0)
		}
	}

	return t
}

type DependencyParameters struct {
	ID         string            `json:"id,omitempty"`
	HbVersion  int               `json:"hbVersion,omitempty"`
	ProjectID  string            `json:"projectId"`
	Shared     bool              `json:"shared,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Parameters map[string]string `json:"parameters"`
}

type Endpoint struct {
	ID           string            `json:"id,omitempty"`
	HbVersion    int               `json:"hbVersion,omitempty"`
	ProjectID    string            `json:"projectId"`
	Shared       bool              `json:"shared,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Type         string            `json:"type"`
	EndpointType string            `json:"endpointType"`
	Endpoint     string            `json:"endpoint"`
	Description  string            `json:"description"`
	Enabled      bool              `json:"enabled"`
	Active       bool              `json:"active"`
}

type HistoryLifecycleEvent struct {
	ID          string            `json:"id,omitempty"`
	HbVersion   int               `json:"hbVersion,omitempty"`
	ProjectID   string            `json:"projectId"`
	Shared      bool              `json:"shared,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Event       string            `json:"event"`
	Description string            `json:"description"`
	ExecutedAt  string            `json:"executedAt"`
}

type Location struct {
	ID        string            `json:"id,omitempty"`
	HbVersion int               `json:"hbVersion,omitempty"`
	ProjectID string            `json:"projectId"`
	Shared    bool              `json:"shared,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Name      string            `json:"name,omitempty"`
	Latitude  string            `json:"latitude,omitempty"`
	Longitude string            `json:"longitude,omitempty"`
}

type BaseNetwork struct {
	ID        string            `json:"id,omitempty"`
	HbVersion int               `json:"hbVersion,omitempty"`
	ProjectID string            `json:"projectId"`
	Shared    bool              `json:"shared,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Name      string            `json:"name"`
	ExtID     string            `json:"extId"`
}

type BaseNfvImage struct {
	ID        string            `json:"id,omitempty"`
	HbVersion int               `json:"hbVersion,omitempty"`
	ProjectID string            `json:"projectId"`
	Shared    bool              `json:"shared,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	ExtID     string            `json:"extId"`
	Created   Date              `json:"created,omitempty"`
}

type Quota struct {
	ID          string            `json:"id,omitempty"`
	HbVersion   int               `json:"hbVersion,omitempty"`
	ProjectID   string            `json:"projectId"`
	Shared      bool              `json:"shared,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Tenant      string            `json:"tenant"`
	Cores       int               `json:"cores"`
	FloatingIPs int               `json:"floatingIps"`
	Instances   int               `json:"instances"`
	KeyPairs    int               `json:"keyPairs"`
	RAM         int               `json:"ram"`
}

type RequiresParameters struct {
	ID         string            `json:"id,omitempty"`
	HbVersion  int               `json:"hbVersion,omitempty"`
	ProjectID  string            `json:"projectId"`
	Shared     bool              `json:"shared,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Parameters []string          `json:"parameters"`
}

type Script struct {
	ID        string            `json:"id,omitempty"`
	HbVersion int               `json:"hbVersion,omitempty"`
	ProjectID string            `json:"projectId"`
	Shared    bool              `json:"shared,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Name      string            `json:"name"`
	Payload   []byte            `json:"-"`
}

type Server struct {
	ID                 string              `json:"id,omitempty"`
	HbVersion          int                 `json:"hbVersion,omitempty"`
	ProjectID          string              `json:"projectId"`
	Shared             bool                `json:"shared,omitempty"`
	Metadata           map[string]string   `json:"metadata,omitempty"`
	Name               string              `json:"name"`
	Image              *DockerImage        `json:"image"`
	Flavour            *DeploymentFlavour  `json:"flavor"`
	Status             string              `json:"status"`
	ExtendedStatus     string              `json:"extendedStatus"`
	ExtID              string              `json:"extId"`
	IPs                map[string][]string `json:"ips"`
	FloatingIPs        map[string]string   `json:"floatingIps"`
	Created            Date                `json:"created,omitempty"`
	Updated            Date                `json:"updated,omitempty"`
	HostName           string              `json:"hostName"`
	HypervisorHostName string              `json:"hypervisorHostName"`
	InstanceName       string              `json:"instanceName"`
}

type Subnet struct {
	ID        string            `json:"id,omitempty"`
	HbVersion int               `json:"hbVersion,omitempty"`
	ProjectID string            `json:"projectId"`
	Shared    bool              `json:"shared,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Name      string            `json:"name"`
	ExtID     string            `json:"extId"`
	NetworkID string            `json:"networkId"`
	CIDR      string            `json:"cidr"`
	GatewayIP string            `json:"gatewayIp"`
}

type BaseVimInstance struct {
	ID        string            `json:"id,omitempty"`
	HbVersion int               `json:"hbVersion,omitempty"`
	ProjectID string            `json:"projectId"`
	Shared    bool              `json:"shared,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Name      string            `json:"name"`
	AuthURL   string            `json:"authUrl"`
	Active    bool              `json:"active"`
	Location  *Location         `json:"location,omitempty"`
	Type      string            `json:"type"`
}

type DockerImage struct {
	BaseNfvImage
	Tags []string `json:"tags"`
}

type DockerNetwork struct {
	BaseNetwork
	Scope   string `json:"scope"`
	Driver  string `json:"driver"`
	Gateway string `json:"gateway"`
	Subnet  string `json:"subnet"`
}

type OpenstackVimInstance struct {
	BaseVimInstance
	Username string `json:"username"`
	Password string `json:"password"`
	KeyPair  string `json:"keyPair"`
	Tenant   string `json:"tenant"`
}

type DockerVimInstance struct {
	BaseVimInstance
	Ca        string          `json:"ca"`
	Cert      string          `json:"cert"`
	DockerKey string          `json:"dockerKey"`
	Images    []DockerImage   `json:"images"`
	Networks  []DockerNetwork `json:"networks"`
}
type VNFCDependencyParameters struct {
	ID         string                           `json:"id,omitempty"`
	HbVersion  int                              `json:"hbVersion,omitempty"`
	ProjectID  string                           `json:"projectId"`
	Shared     bool                             `json:"shared,omitempty"`
	Metadata   map[string]string                `json:"metadata,omitempty"`
	VNFCID     string                           `json:"vnfcId"`
	Parameters map[string]*DependencyParameters `json:"parameters"`
}

type VNFPackage struct {
	ID          string            `json:"id,omitempty"`
	HbVersion   int               `json:"hbVersion,omitempty"`
	ProjectID   string            `json:"projectId"`
	Shared      bool              `json:"shared,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Name        string            `json:"name"`
	NFVOVersion string            `json:"nfvo_version"`
	VIMTypes    []string          `json:"vimTypes"`
	ImageLink   string            `json:"imageLink,omitempty"`
	ScriptsLink string            `json:"scriptsLink,omitempty"`
	Image       *BaseNfvImage     `json:"image,omitempty"`
	Scripts     []*Script         `json:"scripts"`
}

type ManagerCredentials struct {
	RabbitUsername string `json:"rabbitUsername"`
	RabbitPassword string `json:"rabbitPassword"`
}

type PluginRegisterMessage struct {
	Type   string `json:"type"`
	Action string `json:"action"`
}

type VnfmRegisterMessage struct {
	Type     string    `json:"type"`
	Action   string    `json:"action"`
	Endpoint *Endpoint `json:"vnfmManagerEndpoint"`
}

type ManagerUnregisterMessage struct {
	Type     string `json:"type"`
	Action   string `json:"action"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type VnfmManagerUnregisterMessage struct {
	Type     string    `json:"type"`
	Action   string    `json:"action"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Endpoint *Endpoint `json:"vnfmManagerEndpoint"`
}
