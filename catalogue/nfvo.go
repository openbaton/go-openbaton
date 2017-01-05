package catalogue

//go:generate stringer -type=Action
type Action string

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

//go:generate stringer -type=ConfigurationParameter
type ConfigurationParameter struct {
	ID          string `json:"id,omitempty"`
	Version     int    `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	ConfKey     string `json:"confKey"`
	Value       string `json:"value,omitempty"`
}

//go:generate stringer -type=Configuration
type Configuration struct {
	ID                      string                    `json:"id,omitempty"`
	Version                 int                       `json:"version"`
	ProjectID               string                    `json:"projectId"`
	ConfigurationParameters []*ConfigurationParameter `json:"configurationParameters"`
	Name                    string                    `json:"name"`
}

func (cfg *Configuration) Append(p *ConfigurationParameter) {
	cfg.ConfigurationParameters = append(cfg.ConfigurationParameters, p)
}

//go:generate stringer -type=DependencyParameters
type DependencyParameters struct {
	ID         string            `json:"id,omitempty"`
	Version    int               `json:"version"`
	Parameters map[string]string `json:"parameters"`
}

//go:generate stringer -type=Endpoint
type Endpoint struct {
	ID           string `json:"id,omitempty"`
	Version      int    `json:"version"`
	Type         string `json:"type"`
	EndpointType string `json:"endpointType"`
	Endpoint     string `json:"endpoint"`
	Description  string `json:"description"`
	Enabled      bool   `json:"enabled"`
	Active       bool   `json:"active"`
}

//go:generate stringer -type=HistoryLifecycleEvent
type HistoryLifecycleEvent struct {
	ID          string `json:"id,omitempty"`
	Event       string `json:"event"`
	Description string `json:"description"`
	ExecutedAt  string `json:"executedAt"`
}

//go:generate stringer -type=Location
type Location struct {
	ID        string `json:"id,omitempty"`
	Version   int    `json:"id,omitempty"`
	Name      string `json:"id,omitempty"`
	Latitude  string `json:"id,omitempty"`
	Longitude string `json:"id,omitempty"`
}

//go:generate stringer -type=Network
type Network struct {
	ID       string    `json:"id,omitempty"`
	Version  int       `json:"version"`
	Name     string    `json:"name"`
	ExtID    string    `json:"extId"`
	External bool      `json:"external"`
	Shared   bool      `json:"shared"`
	Subnets  []*Subnet `json:"subnets"`
}

//go:generate stringer -type=NFVImage
type NFVImage struct {
	ID              string `json:"id,omitempty"`
	Version         int    `json:"version"`
	ExtID           string `json:"extId"`
	Name            string `json:"name"`
	MinRAM          int64  `json:"minRam"`
	MinDiskSpace    int64  `json:"minDiskSpace"`
	MinCPU          string `json:"minCPU"`
	Public          bool   `json:"public"`
	DiskFormat      string `json:"diskFormat"`
	ContainerFormat string `json:"containerFormat"`
	Created         string `json:"created"` // Actually a date; implement a Marshaller if needed
	Updated         string `json:"updated"` // see above
	IsPublic        bool   `json:"isPublic"`
}

//go:generate stringer -type=RequiresParameters
type RequiresParameters struct {
	ID         string   `json:"id,omitempty"`
	Version    int      `json:"version"`
	Parameters []string `json:"parameters"`
}

//go:generate stringer -type=Script
type Script struct {
	ID      string `json:"id,omitempty"`
	Version int    `json:"version"`
	Name    string `json:"name"`
	Payload []byte `json:"-"`
}

//go:generate stringer -type=Subnet
type Subnet struct {
	ID        string `json:"id,omitempty"`
	Version   int    `json:"version"`
	Name      string `json:"name"`
	ExtID     string `json:"extId"`
	NetworkID string `json:"networkId"`
	CIDR      string `json:"cidr"`
	GatewayIP string `json:"gatewayIp"`
}

//go:generate stringer -type=VIMInstance
type VIMInstance struct {
	ID             string               `json:"id,omitempty"`
	Version        int                  `json:"version"`
	Name           string               `json:"name"`
	AuthURL        string               `json:"authUrl"`
	Tenant         string               `json:"tenant"`
	Username       string               `json:"username"`
	Password       string               `json:"password"`
	KeyPair        string               `json:"keyPair"`
	Location       *Location            `json:"location,omitempty"`
	SecurityGroups []string             `json:"securityGroups"`
	Flavours       []*DeploymentFlavour `json:"flavours"`
	Type           string               `json:"type"`
	Images         []*NFVImage          `json:"images"`
	Networks       []*Network           `json:"networks"`
	ProjectID      string               `json:"projectId"`
	Active         bool                 `json:"active"`
}

func (vi *VIMInstance) HasFlavour(key string) bool {
	for _, df := range vi.Flavours {
		if df.FlavourKey == key || df.ExtID == key || df.ID == key {
			return true
		}
	}

	return false
}

//go:generate stringer -type=VNFCDependencyParameters
type VNFCDependencyParameters struct {
	VNFCID     string                           `json:"vnfcId"`
	ID         string                           `json:"id,omitempty"`
	Version    int                              `json:"version"`
	Parameters map[string]*DependencyParameters `json:"parameters"`
}

//go:generate stringer -type=VNFPackage
type VNFPackage struct {
	ID          string    `json:"id,omitempty"`
	Version     int       `json:"version"`
	Name        string    `json:"name"`
	NFVOVersion string    `json:"nfvo_version"`
	VIMTypes    []string  `json:"vimTypes"`
	ImageLink   string    `json:"imageLink,omitempty"`
	ScriptsLink string    `json:"scriptsLink,omitempty"`
	Image       *NFVImage `json:"image,omitempty"`
	Scripts     []*Script `json:"scripts"`
	ProjectID   string    `json:"projectId"`
}
