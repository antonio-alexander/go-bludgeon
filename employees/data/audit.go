package data

//swagger:model Audit
//Audit is a common type used for accounting of a unique entity (employee)
type Audit struct {
	//The last time (unix nano) something was mutated
	// example: 1652417242000
	LastUpdated int64 `json:"last_updated"`

	//identifies the last someone who mutated something
	// example: bludgeon_employee_memory
	LastUpdatedBy string `json:"last_updated_by"`

	//An integer that's atomically incremented each time something is mutated
	// example: 1
	Version int `json:"version"`
}
