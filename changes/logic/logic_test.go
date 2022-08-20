package logic_test

import (
	"context"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/changes/logic"
	"github.com/antonio-alexander/go-bludgeon/changes/meta"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/file"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/memory"
	"github.com/antonio-alexander/go-bludgeon/changes/meta/mysql"
	"github.com/antonio-alexander/go-bludgeon/internal"

	"github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_meta "github.com/antonio-alexander/go-bludgeon/internal/meta"
	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	mysqlConfig *internal_mysql.Configuration
	fileConfig  *internal_file.Configuration
	logConfig   *logger.Configuration
)

type logicTest struct {
	meta interface {
		meta.Change
		internal.Initializer
	}
	logic.Logic
	internal.Initializer
}

func init() {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if s := strings.Split(env, "="); len(s) > 0 {
			envs[s[0]] = strings.Join(s[1:], ",")
		}
	}
	mysqlConfig = new(internal_mysql.Configuration)
	mysqlConfig.Default()
	mysqlConfig.FromEnv(envs)
	fileConfig = new(internal_file.Configuration)
	fileConfig.Default()
	fileConfig.FromEnv(envs)
	logConfig = new(logger.Configuration)
	logConfig.Default()
	logConfig.FromEnv(envs)
	logConfig.Level = logger.Trace
	logConfig.Prefix = "test_logic"
	rand.Seed(time.Now().UnixNano())
}

func generateId() string {
	return uuid.Must(uuid.NewRandom()).String()
}

func newLogicTest(metaType internal_meta.Type) *logicTest {
	var meta interface {
		meta.Change
		meta.Registration
		meta.RegistrationChange
		internal.Initializer
		internal.Parameterizer
		internal.Configurer
	}

	logger := logger.New()
	logger.Configure(logConfig)
	switch metaType {
	case internal_meta.TypeMemory:
		meta = memory.New()
		meta.SetParameters(logger)
	case internal_meta.TypeFile:
		meta = file.New()
		meta.Configure()
		meta.SetParameters(logger)
		meta.Configure(fileConfig)
	case internal_meta.TypeMySQL:
		meta = mysql.New()
		meta.SetParameters(logger)
		meta.Configure(mysqlConfig)
	}
	logic := logic.New()
	logic.SetParameters(logger, meta)
	return &logicTest{
		meta:        meta,
		Logic:       logic,
		Initializer: logic,
	}
}

func (l *logicTest) initialize(t *testing.T) {
	err := l.Initialize()
	assert.Nil(t, err)
}

func (l *logicTest) shutdown(t *testing.T) {
	l.Shutdown()
	l.meta.Shutdown()
}

func (l *logicTest) testChangeRegistration(t *testing.T) {
	var changes []*data.Change
	ctx := context.TODO()

	//upsert change (neither registration should see this change)
	dataId := generateId()
	dataVersion, dataType := rand.Intn(1000), generateId()
	dataServiceName, whenChanged := generateId(), time.Now().UnixNano()
	changedBy, dataAction := "test_change_crud", generateId()
	changeCreated, err := l.ChangeUpsert(ctx, data.ChangePartial{
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
		l.ChangesDelete(ctx, changeId)
	}(changeCreated.Id)
	changes = append(changes, changeCreated)

	//create registration (1)
	registrationId1 := generateId()
	err = l.RegistrationUpsert(ctx, registrationId1)
	assert.Nil(t, err)
	defer func() {
		l.RegistrationDelete(ctx, registrationId1)
	}()

	//validate that registration (1) doesn't include any changes
	changesRead, err := l.RegistrationChangesRead(ctx, registrationId1)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 0)

	//upsert change (to be seen by registration (1))
	dataId = generateId()
	dataVersion, dataType = rand.Intn(1000), generateId()
	dataServiceName, whenChanged = generateId(), time.Now().UnixNano()
	dataAction, changedBy = generateId(), "test_change_crud"
	changeCreated, err = l.ChangeUpsert(ctx, data.ChangePartial{
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
		l.ChangesDelete(ctx, changeId)
	}(changeCreated.Id)
	changes = append(changes, changeCreated)

	//validate that registration (1) sees the second change, but not the first
	changesRead, err = l.RegistrationChangesRead(ctx, registrationId1)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 1)
	assert.Contains(t, changesRead, changes[1])
	assert.NotContains(t, changesRead, changes[0])

	//create registration
	registrationId2 := generateId()
	err = l.RegistrationUpsert(ctx, registrationId2)
	assert.Nil(t, err)
	defer func() {
		l.RegistrationDelete(ctx, registrationId2)
	}()

	//upsert change (to be seen by registration (1) and (2))
	dataId = generateId()
	dataVersion, dataType = rand.Intn(1000), generateId()
	dataServiceName, whenChanged = generateId(), time.Now().UnixNano()
	dataAction, changedBy = generateId(), "test_change_crud"
	changeCreated, err = l.ChangeUpsert(ctx, data.ChangePartial{
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
		l.ChangesDelete(ctx, changeId)
	}(changeCreated.Id)
	changes = append(changes, changeCreated)

	//validate that registration (2) sees the third change, but not the first or second
	changesRead, err = l.RegistrationChangesRead(ctx, registrationId2)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 1)
	assert.Contains(t, changesRead, changes[2])
	assert.NotContains(t, changesRead, changes[0])
	assert.NotContains(t, changesRead, changes[1])

	//acknowledge the initial change for both services and confirm that there's no change
	err = l.RegistrationChangeAcknowledge(ctx, registrationId1, changes[0].Id)
	assert.Nil(t, err)
	changesRead, err = l.RegistrationChangesRead(ctx, registrationId1)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 2)
	err = l.RegistrationChangeAcknowledge(ctx, registrationId2, changes[0].Id)
	assert.Nil(t, err)
	changesRead, err = l.RegistrationChangesRead(ctx, registrationId2)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 1)

	//validate that initial change has been removed
	change, err := l.ChangeRead(ctx, changes[0].Id)
	assert.NotNil(t, err)
	assert.Nil(t, change)

	//acknowledge the second change for both registrations and confirm actual changes
	err = l.RegistrationChangeAcknowledge(ctx, registrationId1, changes[1].Id)
	assert.Nil(t, err)
	changesRead, err = l.RegistrationChangesRead(ctx, registrationId1)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 1)
	assert.Contains(t, changesRead, changes[2])
	err = l.RegistrationChangeAcknowledge(ctx, registrationId2, changes[1].Id)
	assert.Nil(t, err)
	changesRead, err = l.RegistrationChangesRead(ctx, registrationId2)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 1)
	assert.Contains(t, changesRead, changes[2])

	//validate that second change has been removed
	change, err = l.ChangeRead(ctx, changes[1].Id)
	assert.NotNil(t, err)
	assert.Nil(t, change)

	//delete the initial registration, then re-create it and ensure that there are no changes
	err = l.RegistrationDelete(ctx, registrationId1)
	assert.Nil(t, err)
	err = l.RegistrationUpsert(ctx, registrationId1)
	assert.Nil(t, err)
	changesRead, err = l.RegistrationChangesRead(ctx, registrationId1)
	assert.Nil(t, err)
	assert.Len(t, changesRead, 0)
}

func (l *logicTest) testChangeHandlers(t *testing.T) {
	var changeId string

	ctx := context.TODO()
	changeReceived := make(chan struct{})

	//create handler
	handlerId, err := l.HandlerCreate(ctx, func(ctx context.Context, handlerId string, changes []*data.Change) error {
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
		l.HandlerDelete(ctx, handlerId)
	}()

	//upsert change
	dataId := generateId()
	dataVersion, dataType := rand.Intn(1000), generateId()
	dataServiceName, whenChanged := generateId(), time.Now().UnixNano()
	dataAction, changedBy := generateId(), "test_change_crud"
	changeCreated, err := l.ChangeUpsert(ctx, data.ChangePartial{
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
		l.ChangesDelete(ctx, changeId)
	}(changeCreated.Id)
	changeId = changeCreated.Id

	//validate change received
	select {
	case <-changeReceived:
	case <-time.After(10 * time.Second):
		assert.Fail(t, "unable to confirm change received")
	}
}

func testLogic(t *testing.T, metaType internal_meta.Type) {
	l := newLogicTest(internal_meta.TypeMemory)

	//initialize
	l.initialize(t)
	defer l.shutdown(t)

	//execute tests
	t.Run("Change Registration", l.testChangeRegistration)
	t.Run("Change Handlers", l.testChangeHandlers)
}

func TestLogicMemory(t *testing.T) {
	testLogic(t, internal_meta.TypeMemory)
}

func TestLogicFile(t *testing.T) {
	testLogic(t, internal_meta.TypeFile)
}

func TestLogicMysql(t *testing.T) {
	testLogic(t, internal_meta.TypeMySQL)
}
