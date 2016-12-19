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
)

type ConfigurationParameter struct {
	ID          string `json:"id"`
	Version     int    `json:"version"`
	Description string `json:"description"`
	ConfKey     string `json:"confKey"`
	Value       string `json:"value"`
}

type Configuration struct {
	ID                      string                    `json:"id"`
	Version                 int                       `json:"version"`
	ProjectID               string                    `json:"projectId"`
	ConfigurationParameters []*ConfigurationParameter `json:"configurationParameters"`
	Name                    string                    `json:"name"`
}

type DependencyParameters struct {
	ID         string            `json:"id"`
	Version    int               `json:"version"`
	Parameters map[string]string `json:"parameters"`
}

type HistoryLifecycleEvent struct {
	ID          string `json:"id"`
	Event       string `json:"event"`
	Description string `json:"description"`
	ExecutedAt  string `json:"executedAt"`
}

type Location struct {
	ID string `json:"id"`
	Version int `json:"id"`
	Name string `json:"id"`
	Latitude string `json:"id"`
	Longitude string `json:"id"`
}

type Script struct {
	ID string `json:"id"`
	Version int `json:"version"`
	Name string `json:"name"`
	Payload []byte `json:"-"`
}

type VIMInstance struct {
	ID string `json:"id"`
	Version int `json:"version"`
	Name string `json:"name"`
	AuthURL string `json:"authUrl"`
	Tenant string `json:"tenant"`
	Username string `json:"username"`
	Password string `json:"password"`
	KeyPair string `json:"keyPair"`
	Location *Location `json:"location"`
	SecurityGroups []string `json:"securityGroups"`
	Flavours []*DeploymentFlavour `json:"flavours"`
	Type string `json:"type"`
	Images []*NFVImage `json:"images"`
	Networks []*Network `json:"networks"`
	ProjectID string `json:"projectId"`
	Active bool `json:"active"`
}

type VNFCDependencyParameters struct {
	VNFCID     string                           `json:"vnfcId"`
	ID         string                           `json:"id"`
	Version    int                              `json:"version"`
	Parameters map[string]*DependencyParameters `json:"parameters"`
}
