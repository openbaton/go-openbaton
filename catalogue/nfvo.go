package catalogue

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

type ConfigurationParameter struct {
	ID          ID     `json:"id"`
	Version     int    `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	ConfKey     string `json:"confKey"`
	Value       string `json:"value,omitempty"`
}

type Configuration struct {
	ID                      ID                        `json:"id"`
	Version                 int                       `json:"version"`
	ProjectID               string                    `json:"projectId"`
	ConfigurationParameters []*ConfigurationParameter `json:"configurationParameters"`
	Name                    string                    `json:"name"`
}

func (cfg *Configuration) Append(p *ConfigurationParameter) {
	cfg.ConfigurationParameters = append(cfg.ConfigurationParameters, p)
}

type DependencyParameters struct {
	ID         ID                `json:"id"`
	Version    int               `json:"version"`
	Parameters map[string]string `json:"parameters"`
}

type HistoryLifecycleEvent struct {
	ID          ID     `json:"id"`
	Event       string `json:"event"`
	Description string `json:"description"`
	ExecutedAt  string `json:"executedAt"`
}

type Location struct {
	ID        ID     `json:"id"`
	Version   int    `json:"id"`
	Name      string `json:"id"`
	Latitude  string `json:"id"`
	Longitude string `json:"id"`
}

type Network struct {
	ID       ID        `json:"id"`
	Version  int       `json:"version"`
	Name     string    `json:"name"`
	ExtID    string    `json:"extId"`
	External bool      `json:"external"`
	Shared   bool      `json:"shared"`
	Subnets  []*Subnet `json:"subnets"`
}

type NFVImage struct {
	ID              ID     `json:"id"`
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

type RequiresParameters struct {
	ID         ID       `json:"id"`
	Version    int      `json:"version"`
	Parameters []string `json:"parameters"`
}

type Script struct {
	ID      ID     `json:"id"`
	Version int    `json:"version"`
	Name    string `json:"name"`
	Payload []byte `json:"-"`
}

type Subnet struct {
	ID        ID     `json:"id"`
	Version   int    `json:"version"`
	Name      string `json:"name"`
	ExtID     string `json:"extId"`
	NetworkID string `json:"networkId"`
	CIDR      string `json:"cidr"`
	GatewayIP string `json:"gatewayIp"`
}

type VIMInstance struct {
	ID             ID                   `json:"id"`
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

type VNFCDependencyParameters struct {
	VNFCID     string                           `json:"vnfcId"`
	ID         ID                               `json:"id"`
	Version    int                              `json:"version"`
	Parameters map[string]*DependencyParameters `json:"parameters"`
}

type VNFPackage struct {
	ID          ID        `json:"id"`
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
