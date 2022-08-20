package memory

import (
	data "github.com/antonio-alexander/go-bludgeon/changes/data"

	"github.com/google/uuid"
)

func generateID() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

func copyChange(c *data.Change) *data.Change {
	return &data.Change{
		Id:              c.Id,
		WhenChanged:     c.WhenChanged,
		ChangedBy:       c.ChangedBy,
		DataId:          c.DataId,
		DataServiceName: c.DataServiceName,
		DataType:        c.DataType,
		DataAction:      c.DataAction,
		DataVersion:     c.DataVersion,
	}
}
