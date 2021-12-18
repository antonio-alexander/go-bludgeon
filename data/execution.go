package data

import (
	uuid "github.com/google/uuid"
)

func GenerateID() (string, error) {
	guid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return guid.String(), nil
}
