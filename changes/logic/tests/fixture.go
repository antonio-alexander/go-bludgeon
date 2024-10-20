package tests

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/changes/logic"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type Fixture struct {
	logic.Logic
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateId() string {
	return uuid.Must(uuid.NewRandom()).String()
}

func NewFixture(items ...interface{}) *Fixture {
	f := &Fixture{}
	for _, item := range items {
		switch item := item.(type) {
		case logic.Logic:
			f.Logic = item
		}
	}
	return f
}

func (f *Fixture) TestChangeRegistration(t *testing.T) {
	var changes []*data.Change
	ctx := context.TODO()

	//upsert change (neither registration should see this change)
	dataId := generateId()
	dataVersion, dataType := rand.Intn(1000), generateId()
	dataServiceName, whenChanged := generateId(), time.Now().UnixNano()
	changedBy, dataAction := "test_change_crud", generateId()
	changeCreated, err := f.ChangeUpsert(ctx, data.ChangePartial{
		DataId:          &dataId,
		DataVersion:     &dataVersion,
		DataType:        &dataType,
		DataServiceName: &dataServiceName,
		DataAction:      &dataAction,
		WhenChanged:     &whenChanged,
		ChangedBy:       &changedBy,
	})
	assert.Nil(t, err)
	assert.NotNil(t, changeCreated)
	defer func(changeId string) {
		f.ChangesDelete(ctx, changeId)
	}(changeCreated.Id)
	changes = append(changes, changeCreated)

	//create registration (1)
	registrationId1 := generateId()
	err = f.RegistrationUpsert(ctx, registrationId1)
	assert.Nil(t, err)
	defer func() {
		f.RegistrationDelete(ctx, registrationId1)
	}()

	//validate that registration (1) doesn't include any changes
	changesRead, err := f.RegistrationChangesRead(ctx, registrationId1)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 0)

	//upsert change (to be seen by registration (1))
	dataId = generateId()
	dataVersion, dataType = rand.Intn(1000), generateId()
	dataServiceName, whenChanged = generateId(), time.Now().UnixNano()
	dataAction, changedBy = generateId(), "test_change_crud"
	changeCreated, err = f.ChangeUpsert(ctx, data.ChangePartial{
		DataId:          &dataId,
		DataVersion:     &dataVersion,
		DataType:        &dataType,
		DataServiceName: &dataServiceName,
		DataAction:      &dataAction,
		WhenChanged:     &whenChanged,
		ChangedBy:       &changedBy,
	})
	assert.Nil(t, err)
	assert.NotNil(t, changeCreated)
	defer func(changeId string) {
		f.ChangesDelete(ctx, changeId)
	}(changeCreated.Id)
	changes = append(changes, changeCreated)

	//validate that registration (1) sees the second change, but not the first
	changesRead, err = f.RegistrationChangesRead(ctx, registrationId1)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 1)
	assert.Contains(t, changesRead, changes[1])
	assert.NotContains(t, changesRead, changes[0])

	//create registration
	registrationId2 := generateId()
	err = f.RegistrationUpsert(ctx, registrationId2)
	assert.Nil(t, err)
	defer func() {
		f.RegistrationDelete(ctx, registrationId2)
	}()

	//upsert change (to be seen by registration (1) and (2))
	dataId = generateId()
	dataVersion, dataType = rand.Intn(1000), generateId()
	dataServiceName, whenChanged = generateId(), time.Now().UnixNano()
	dataAction, changedBy = generateId(), "test_change_crud"
	changeCreated, err = f.ChangeUpsert(ctx, data.ChangePartial{
		DataId:          &dataId,
		DataVersion:     &dataVersion,
		DataType:        &dataType,
		DataServiceName: &dataServiceName,
		DataAction:      &dataAction,
		WhenChanged:     &whenChanged,
		ChangedBy:       &changedBy,
	})
	assert.Nil(t, err)
	assert.NotNil(t, changeCreated)
	defer func(changeId string) {
		f.ChangesDelete(ctx, changeId)
	}(changeCreated.Id)
	changes = append(changes, changeCreated)

	//validate that registration (2) sees the third change, but not the first or second
	changesRead, err = f.RegistrationChangesRead(ctx, registrationId2)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 1)
	assert.Contains(t, changesRead, changes[2])
	assert.NotContains(t, changesRead, changes[0])
	assert.NotContains(t, changesRead, changes[1])

	//acknowledge the initial change for both services and confirm that there's no change
	err = f.RegistrationChangeAcknowledge(ctx, registrationId1, changes[0].Id)
	assert.Nil(t, err)
	changesRead, err = f.RegistrationChangesRead(ctx, registrationId1)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 2)
	err = f.RegistrationChangeAcknowledge(ctx, registrationId2, changes[0].Id)
	assert.Nil(t, err)
	changesRead, err = f.RegistrationChangesRead(ctx, registrationId2)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 1)

	//validate that initial change has been removed
	change, err := f.ChangeRead(ctx, changes[0].Id)
	assert.NotNil(t, err)
	assert.Nil(t, change)

	//acknowledge the second change for both registrations and confirm actual changes
	err = f.RegistrationChangeAcknowledge(ctx, registrationId1, changes[1].Id)
	assert.Nil(t, err)
	changesRead, err = f.RegistrationChangesRead(ctx, registrationId1)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 1)
	assert.Contains(t, changesRead, changes[2])
	err = f.RegistrationChangeAcknowledge(ctx, registrationId2, changes[1].Id)
	assert.Nil(t, err)
	changesRead, err = f.RegistrationChangesRead(ctx, registrationId2)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 1)
	assert.Contains(t, changesRead, changes[2])

	//validate that second change has been removed
	change, err = f.ChangeRead(ctx, changes[1].Id)
	assert.NotNil(t, err)
	assert.Nil(t, change)

	//delete the initial registration, then re-create it and ensure that there are no changes
	err = f.RegistrationDelete(ctx, registrationId1)
	assert.Nil(t, err)
	err = f.RegistrationUpsert(ctx, registrationId1)
	assert.Nil(t, err)
	changesRead, err = f.RegistrationChangesRead(ctx, registrationId1)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 0)
}

func (f *Fixture) TestChangeHandlers(t *testing.T) {
	var changeId string

	ctx := context.TODO()
	changeReceived := make(chan struct{})

	//create handler
	handlerId, err := f.HandlerCreate(ctx, func(ctx context.Context, handlerId string, changes []*data.Change) error {
		for _, change := range changes {
			if changeId == change.Id {
				select {
				default:
					close(changeReceived)
				case <-changeReceived:
				}
			}
		}
		return nil
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, handlerId)
	defer func() {
		f.HandlerDelete(ctx, handlerId)
	}()

	//upsert change
	dataId := generateId()
	dataVersion, dataType := rand.Intn(1000), generateId()
	dataServiceName, whenChanged := generateId(), time.Now().UnixNano()
	dataAction, changedBy := generateId(), "test_change_crud"
	changeCreated, err := f.ChangeUpsert(ctx, data.ChangePartial{
		DataId:          &dataId,
		DataVersion:     &dataVersion,
		DataType:        &dataType,
		DataServiceName: &dataServiceName,
		DataAction:      &dataAction,
		WhenChanged:     &whenChanged,
		ChangedBy:       &changedBy,
	})
	assert.Nil(t, err)
	assert.NotNil(t, changeCreated)
	defer func(changeId string) {
		f.ChangesDelete(ctx, changeId)
	}(changeCreated.Id)
	changeId = changeCreated.Id

	//validate change received
	select {
	case <-changeReceived:
	case <-time.After(10 * time.Second):
		assert.Fail(t, "unable to confirm change received")
	}
}
