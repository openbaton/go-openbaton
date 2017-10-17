package catalogue

type BaseEntity struct {
	ID        string            `json:"id,omitempty"`
	HbVersion int               `json:"hbVersion,omitempty"`
	ProjectID string            `json:"projectId"`
	Shared    bool              `json:"shared,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}
