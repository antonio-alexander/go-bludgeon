package tests

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/changes/meta"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateId() string {
	return uuid.Must(uuid.NewRandom()).String()
}

func TestChangeCRUD(m interface {
	meta.Change
}) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.TODO()

		//create change
		dataId := generateId()
		dataVersion, dataType := rand.Intn(1000), "employee"
		dataServiceName, whenChanged := "employees", time.Now().UnixNano()
		changedBy, dataAction := "test_change_crud", "create"
		changeCreated, err := m.ChangeCreate(ctx, data.ChangePartial{
			DataId:          &dataId,
			DataVersion:     &dataVersion,
			DataType:        &dataType,
			DataAction:      &dataAction,
			DataServiceName: &dataServiceName,
			WhenChanged:     &whenChanged,
			ChangedBy:       &changedBy,
		})
		assert.Nil(t, err)
		assert.Equal(t, dataId, changeCreated.DataId)
		assert.Equal(t, dataVersion, changeCreated.DataVersion)
		assert.Equal(t, dataType, changeCreated.DataType)
		assert.Equal(t, dataServiceName, changeCreated.DataServiceName)
		assert.Equal(t, dataAction, changeCreated.DataAction)
		//TODO: fix this
		// assert.Equal(t, whenChanged, changeCreated.WhenChanged)
		assert.Equal(t, changedBy, changeCreated.ChangedBy)
		changeId := changeCreated.Id
		defer func() {
			m.ChangesDelete(ctx, changeId)
		}()

		//read change
		changeRead, err := m.ChangeRead(ctx, changeId)
		assert.Nil(t, err)
		assert.Equal(t, changeCreated, changeRead)

		//delete change
		err = m.ChangesDelete(ctx, changeId)
		assert.Nil(t, err)

		//read change
		changeRead, err = m.ChangeRead(ctx, changeId)
		assert.NotNil(t, err)
		assert.Nil(t, changeRead)
	}
}

func TestChangeSearch(m interface {
	meta.Change
}) func(*testing.T) {
	return func(t *testing.T) {
		var changesCreated []*data.Change
		var changeIds, dataIds, dataTypes []string
		var dataServiceNames []string
		var dataVersions []int

		ctx := context.TODO()

		//generate dynamic constants
		for i := 0; i < 5; i++ {
			dataIds = append(dataIds, generateId())
			dataVersions = append(dataVersions, rand.Intn(1000))
		}
		for i := 0; i < 3; i++ {
			dataTypes = append(dataTypes, generateId())
			dataServiceNames = append(dataServiceNames, generateId())
		}
		whenChanged, changedBy := time.Now().UnixNano(), "test_change_crud"

		//create changes
		for _, changePartial := range []data.ChangePartial{
			{
				DataId:          &dataIds[0],
				DataVersion:     &dataVersions[0],
				DataType:        &dataTypes[0],
				DataServiceName: &dataServiceNames[0],
				WhenChanged:     &whenChanged,
				ChangedBy:       &changedBy,
			},
			{
				DataId:          &dataIds[1],
				DataVersion:     &dataVersions[1],
				DataType:        &dataTypes[1],
				DataServiceName: &dataServiceNames[1],
				WhenChanged:     &whenChanged,
				ChangedBy:       &changedBy,
			},
			{
				DataId:          &dataIds[2],
				DataVersion:     &dataVersions[2],
				DataType:        &dataTypes[2],
				DataServiceName: &dataServiceNames[2],
				WhenChanged:     &whenChanged,
				ChangedBy:       &changedBy,
			},
			{
				DataId:          &dataIds[3],
				DataVersion:     &dataVersions[3],
				DataType:        &dataTypes[0],
				DataServiceName: &dataServiceNames[0],
				WhenChanged:     &whenChanged,
				ChangedBy:       &changedBy,
			},
			{
				DataId:          &dataIds[4],
				DataVersion:     &dataVersions[4],
				DataType:        &dataTypes[1],
				DataServiceName: &dataServiceNames[1],
				WhenChanged:     &whenChanged,
				ChangedBy:       &changedBy,
			},
		} {
			changeCreated, err := m.ChangeCreate(ctx, changePartial)
			assert.Nil(t, err)
			assert.NotNil(t, changeCreated)
			changesCreated = append(changesCreated, changeCreated)
			changeIds = append(changeIds, changeCreated.Id)
		}
		defer func() {
			m.ChangesDelete(ctx, changeIds...)
		}()

		//attempt to read changes by change id
		changesRead, err := m.ChangesRead(ctx, data.ChangeSearch{
			ChangeIds: changeIds,
		})
		assert.Nil(t, err)
		assert.Len(t, changesRead, len(changesCreated))
		assert.Contains(t, changesRead, changesCreated[0])
		assert.Contains(t, changesRead, changesCreated[1])
		assert.Contains(t, changesRead, changesCreated[2])
		assert.Contains(t, changesRead, changesCreated[3])
		assert.Contains(t, changesRead, changesCreated[4])

		//attempt to read changes by data id
		changesRead, err = m.ChangesRead(ctx, data.ChangeSearch{
			DataIds: dataIds,
		})
		assert.Nil(t, err)
		assert.Len(t, changesRead, len(changesCreated))
		assert.Contains(t, changesRead, changesCreated[0])
		assert.Contains(t, changesRead, changesCreated[1])
		assert.Contains(t, changesRead, changesCreated[2])
		assert.Contains(t, changesRead, changesCreated[3])
		assert.Contains(t, changesRead, changesCreated[4])

		//TODO: types
		//TODO: services
		//TODO: since
		//TODO: version
	}
}

func TestRegistrationCRUD(m interface {
	meta.Registration
}) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.TODO()

		//upsert registration
		registrationId := generateId()
		err := m.RegistrationUpsert(ctx, registrationId)
		assert.Nil(t, err)

		//delete registration
		err = m.RegistrationDelete(ctx, registrationId)
		assert.Nil(t, err)

		//delete registration
		err = m.RegistrationDelete(ctx, registrationId)
		assert.NotNil(t, err)
	}
}

func TestRegistrationChanges(m interface {
	meta.Change
	meta.Registration
	meta.RegistrationChange
}) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.TODO()

		//upsert registration
		registrationId := generateId()
		err := m.RegistrationUpsert(ctx, registrationId)
		assert.Nil(t, err)
		defer func() {
			m.RegistrationDelete(ctx, registrationId)
		}()

		//upsert change
		dataId := generateId()
		dataVersion, dataType := rand.Intn(1000), generateId()
		dataService, whenChanged := generateId(), time.Now().UnixNano()
		changedBy := "test_change_crud"
		changeCreated, err := m.ChangeCreate(ctx, data.ChangePartial{
			DataId:          &dataId,
			DataVersion:     &dataVersion,
			DataType:        &dataType,
			DataServiceName: &dataService,
			WhenChanged:     &whenChanged,
			ChangedBy:       &changedBy,
		})
		assert.Nil(t, err)
		assert.Equal(t, dataId, changeCreated.DataId)
		assert.Equal(t, dataVersion, changeCreated.DataVersion)
		assert.Equal(t, dataType, changeCreated.DataType)
		assert.Equal(t, dataService, changeCreated.DataServiceName)
		//TODO: fix this
		// assert.Equal(t, whenChanged, changeCreated.WhenChanged)
		assert.Equal(t, changedBy, changeCreated.ChangedBy)
		changeId := changeCreated.Id
		defer func() {
			m.ChangesDelete(ctx, changeId)
		}()

		//upsert registration changes
		err = m.RegistrationChangeUpsert(ctx, changeId)
		assert.Nil(t, err)

		//read registration changes
		changesRead, err := m.RegistrationChangesRead(ctx, registrationId)
		assert.Nil(t, err)
		assert.Len(t, changesRead, 1)
		assert.Contains(t, changesRead, changeId)

		//attempt to delete change (should fail)
		err = m.ChangesDelete(ctx, changeId)
		assert.NotNil(t, err)

		//acknowledge change
		changeIdsToPrune, err := m.RegistrationChangeAcknowledge(ctx, registrationId, changeId)
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, len(changeIdsToPrune), 1)
		assert.Contains(t, changeIdsToPrune, changeId)

		//delete change
		err = m.ChangesDelete(ctx, changeId)
		assert.Nil(t, err)

		//read registration changes
		changesRead, err = m.RegistrationChangesRead(ctx, registrationId)
		assert.Nil(t, err)
		assert.Len(t, changesRead, 0)
		assert.NotContains(t, changesRead, changeId)

		//delete registration
		err = m.RegistrationDelete(ctx, registrationId)
		assert.Nil(t, err)
	}
}
