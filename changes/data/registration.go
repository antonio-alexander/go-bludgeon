package data

import "encoding/json"

// swagger:model Registration
type Registration struct {
	// The ID of the unique change that has occurred (v4 UUID)
	// example: 86fa2f09-d260-11ec-bd5d-0242c0a8e002
	Id string `json:"id"`

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

type RegistrationPartial struct {
	// The time the change occurred
	// example: 1652417242000
	WhenRegistrationd *int64 `json:"when_changed,string,omitempty"`

	// Identifies the someone that performed the change
	// example: bludgeon_employee_memory
	RegistrationdBy *string `json:"changed_by,omitempty"`

	// The ID of the underlying data that has been changed (v4 UUID)
	// example: 86fa2f09-d260-11ec-bd5d-0242c0a8e002
	DataId *string `json:"data_id,omitempty"`

	// A string that identifies the service that the change belongs to
	// example: employees
	DataServiceName *string `json:"data_service_name,omitempty"`

	// A string that identifies the data type that was changed
	// example: employee
	DataType *string `json:"data_type,omitempty"`

	//A string that identifies the action that has occured to the data
	// example: create
	DataAction *string `json:"data_action,omitempty"`

	// An integer that's atomically incremented each time something is mutated
	// example: 1
	DataVersion *int `json:"data_version,omitempty"`
}

func (c *RegistrationPartial) Type() MessageType {
	return MessageTypeRegistrationPartial
}

func (c *RegistrationPartial) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *RegistrationPartial) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, c)
}

type RegistrationDigest struct {
	Registrations []*Registration `json:"registrations"`
}

func (r *RegistrationDigest) Type() MessageType {
	return MessageTypeRegistrationDigest
}

func (r *RegistrationDigest) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *RegistrationDigest) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, r)
}
