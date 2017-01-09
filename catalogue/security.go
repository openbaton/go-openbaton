package catalogue

type Key struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	ProjectID   string `json:"projectId"`
	PublicKey   string `json:"publicKey"`
	Fingerprint string `json:"fingerprint"`
}
