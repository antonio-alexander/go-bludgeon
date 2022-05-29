package memory

import (
	data "github.com/antonio-alexander/go-bludgeon/employees/data"

	"github.com/google/uuid"
)

func generateID() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

func copyEmployee(e *data.Employee) *data.Employee {
	return &data.Employee{
		ID:           e.ID,
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		EmailAddress: e.EmailAddress,
		Audit: data.Audit{
			LastUpdated:   e.LastUpdated,
			LastUpdatedBy: e.LastUpdatedBy,
			Version:       e.Version,
		},
	}
}
