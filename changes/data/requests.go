package data

import (
	"encoding/json"
)

type RequestChange struct {
	ChangePartial ChangePartial `json:"change_partial"`
}

func (r *RequestChange) Type() MessageType {
	return MessageTypeRequestChange
}

func (r *RequestChange) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *RequestChange) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, r)
}

type RequestAcknowledge struct {
	ChangeIds []string `json:"change_ids"`
}

func (r *RequestAcknowledge) Type() MessageType {
	return MessageTypeRequestAcknowledge
}

func (r *RequestAcknowledge) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *RequestAcknowledge) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, r)
}

type RequestRegister struct {
	RegistrationId string `json:"registration_id"`
}

func (r *RequestRegister) Type() MessageType {
	return MessageTypeRequestRegister
}

func (r *RequestRegister) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *RequestRegister) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, r)
}
