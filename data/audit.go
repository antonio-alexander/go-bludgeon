package data

//Audit is a common type we use to represent the fields
// that are auditable
type Audit struct {
	LastUpdated   int64  `json:"last_updated"`
	LastUpdatedBy string `json:"last_updated_by"`
	Version       int    `json:"version"`
}
