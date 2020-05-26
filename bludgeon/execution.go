package bludgeon

import (
	uuid "github.com/google/uuid"
)

func GenerateID() (id string, err error) {
	var guid uuid.UUID

	//create uuid
	if guid, err = uuid.NewRandom(); err != nil {
		return
	}
	//get id in string format
	id = guid.String()

	return
}
