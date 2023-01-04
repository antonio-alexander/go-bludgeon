package data

import (
	"encoding/json"
)

type ResponseAcknowledge struct {
	RegistrationId string   `json:"registration_id"`
	ChangeIds      []string `json:"change_ids"`
}

func (r *ResponseAcknowledge) Type() MessageType {
	return MessageTypeResponseAcknowledge
}

func (r *ResponseAcknowledge) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ResponseAcknowledge) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, r)
}

type ResponseRegister struct {
	RegistrationId string `json:"registration_id"`
}

func (r *ResponseRegister) Type() MessageType {
	return MessageTypeResponseRegister
}

func (r *ResponseRegister) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ResponseRegister) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, r)
}
