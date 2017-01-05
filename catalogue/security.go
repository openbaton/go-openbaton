package catalogue

type Key struct {
	ID          ID     `json:"id,omitempty"`
	Name        string `json:"name"`
	ProjectID   string `json:"projectId"`
	PublicKey   string `json:"publicKey"`
	Fingerprint string `json:"fingerprint"`
}
