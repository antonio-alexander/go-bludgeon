package data

import "encoding/json"

type Change struct {
	// The ID of the unique change that has occurred (v4 UUID)
	// example: 86fa2f09-d260-11ec-bd5d-0242c0a8e002
	Id string `json:"id"`

	// The time the change occurred
	// example: 1652417242000
	WhenChanged int64 `json:"when_changed,string"`

	// Identifies the someone that performed the change
	// example: bludgeon_employee_memory
	ChangedBy string `json:"changed_by"`

	// The ID of the underlying data that has been changed (v4 UUID)
	// example: 86fa2f09-d260-11ec-bd5d-0242c0a8e002
	DataId string `json:"data_id"`

	// A string that identifies the service that the change belongs to
	// example: employees
	DataServiceName string `json:"data_service_name"`

	// A string that identifies the data type that was changed
	// example: employee
	DataType string `json:"data_type"`

	//A string that identifies the action that has occured to the data
	// example: create
	DataAction string `json:"data_action"`

	// An integer that's atomically incremented each time something is mutated
	// example: 1
	DataVersion int `json:"data_version"`
}

func (c *Change) Type() MessageType {
	return MessageTypeChange
}

func (c *Change) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Change) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, c)
}

type ChangeDigest struct {
	Changes []*Change `json:"changes"`
}

func (c *ChangeDigest) Type() MessageType {
	return MessageTypeChangeDigest
}

func (c *ChangeDigest) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *ChangeDigest) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, c)
}

type ChangePartial struct {
	// The time the change occurred
	// example: 1652417242000
	WhenChanged *int64 `json:"when_changed,string,omitempty"`

	// Identifies the someone that performed the change
	// example: bludgeon_employee_memory
	ChangedBy *string `json:"changed_by,omitempty"`

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
