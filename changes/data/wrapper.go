package data

import (
	"encoding"
	"encoding/json"

	"github.com/pkg/errors"
)

type MessageType string

const (
	MessageTypeEmpty               MessageType = "empty"
	MessageTypeChange              MessageType = "change"
	MessageTypeChangePartial       MessageType = "change_partial"
	MessageTypeChangeDigest        MessageType = "change_digest"
	MessageTypeRequestRegister     MessageType = "request_register"
	MessageTypeRequestAcknowledge  MessageType = "request_acknowledge"
	MessageTypeResponseRegister    MessageType = "response_register"
	MessageTypeResponseAcknowledge MessageType = "response_acknowledge"
	MessageTypeRegistration        MessageType = "registration"
	MessageTypeRegistrationPartial MessageType = "registration_partial"
	MessageTypeRegistrationDigest  MessageType = "registration_digest"
)

type Empty struct{}

func (e *Empty) Type() MessageType {
	return MessageTypeEmpty
}

func (e *Empty) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

type Wrapper struct {
	Type  MessageType `json:"type"`
	Bytes []byte      `json:"bytes"`
}

func (w *Wrapper) MarshalBinary() ([]byte, error) {
	return json.Marshal(w)
}

type Wrappable interface {
	Type() MessageType
	encoding.BinaryMarshaler
}

func ToWrapper(message Wrappable) *Wrapper {
	bytes, _ := message.MarshalBinary()
	return &Wrapper{
		Type:  message.Type(),
		Bytes: bytes,
	}
}

func FromWrapper(wrapper *Wrapper) (interface{}, error) {
	if wrapper == nil {
		return nil, errors.New("wrapper is nil")
	}
	switch wrapper.Type {
	case MessageTypeEmpty:
		return nil, nil
	case MessageTypeChange:
		change := &Change{}
		if err := change.UnmarshalBinary(wrapper.Bytes); err != nil {
			return nil, err
		}
		return change, nil
	case MessageTypeChangeDigest:
		changeDigest := &ChangeDigest{}
		if err := changeDigest.UnmarshalBinary(wrapper.Bytes); err != nil {
			return nil, err
		}
		return changeDigest, nil
	case MessageTypeRequestAcknowledge:
		request := &RequestAcknowledge{}
		if err := request.UnmarshalBinary(wrapper.Bytes); err != nil {
			return nil, err
		}
		return request, nil
	case MessageTypeRequestRegister:
		request := &RequestRegister{}
		if err := request.UnmarshalBinary(wrapper.Bytes); err != nil {
			return nil, err
		}
		return request, nil
	case MessageTypeResponseRegister:
		response := &ResponseRegister{}
		if err := response.UnmarshalBinary(wrapper.Bytes); err != nil {
			return nil, err
		}
		return response, nil
	case MessageTypeResponseAcknowledge:
		response := &ResponseAcknowledge{}
		if err := response.UnmarshalBinary(wrapper.Bytes); err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, errors.Errorf("unsupported type: %s", wrapper.Type)
}
