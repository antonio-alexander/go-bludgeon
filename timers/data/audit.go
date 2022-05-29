package data

//Audit is a common type used for accounting of a unique entity (employee)
// swagger:model Audit
type Audit struct {
	//LastUpdated represents the last time (unix nano) something was mutated
	// example: 1652417242000
	LastUpdated int64 `json:"last_updated"`

	//LastUpdatedBy will identify the last someone who mutated something
	// example: bludgeon_employee_memory
	LastUpdatedBy string `json:"last_updated_by"`

	//Version is an integer that's atomically incremented each time something i smutated
	// example: 1
	Version int `json:"version"`
}
