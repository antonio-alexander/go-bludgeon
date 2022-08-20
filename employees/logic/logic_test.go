package logic_test

import (
	"context"
	"math/rand"
	"os"
	"path"
	"strings"
	"testing"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	logic "github.com/antonio-alexander/go-bludgeon/employees/logic"
	meta "github.com/antonio-alexander/go-bludgeon/employees/meta"
	file "github.com/antonio-alexander/go-bludgeon/employees/meta/file"
	memory "github.com/antonio-alexander/go-bludgeon/employees/meta/memory"
	mysql "github.com/antonio-alexander/go-bludgeon/employees/meta/mysql"

	changesclient "github.com/antonio-alexander/go-bludgeon/changes/client"
	changesclientkafka "github.com/antonio-alexander/go-bludgeon/changes/client/kafka"
	changesclientrest "github.com/antonio-alexander/go-bludgeon/changes/client/rest"
	changesdata "github.com/antonio-alexander/go-bludgeon/changes/data"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	internal_kafka "github.com/antonio-alexander/go-bludgeon/internal/kafka"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	"github.com/stretchr/testify/assert"
)

const filename string = "bludgeon_logic.json"

var (
	configMetaMysql          = new(internal_mysql.Configuration)
	configMetaFile           = new(internal_file.Configuration)
	configChangesClientRest  = new(changesclientrest.Configuration)
	configChangesClientKafka = new(changesclientkafka.Configuration)
	configKafkaClient        = new(internal_kafka.Configuration)
	letterRunes              = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type logicTest struct {
	meta interface {
		internal.Initializer
		internal.Configurer
		meta.Employee
	}
	changesClient interface {
		internal.Initializer
		internal.Configurer
		changesclient.Client
	}
	changesHandler interface {
		internal.Initializer
		internal.Configurer
		changesclient.Handler
	}
	logic.Logic
}

func init() {
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		if s := strings.Split(e, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	configMetaMysql.Default()
	configMetaMysql.FromEnv(envs)
	configMetaFile.Default()
	configMetaFile.FromEnv(envs)
	configMetaFile.File = path.Join("../../tmp", filename)
	os.Remove(configMetaFile.File)
	configChangesClientRest.Default()
	configChangesClientRest.FromEnv(envs)
	configKafkaClient.Default()
	configKafkaClient.Brokers = []string{"localhost:9092"}
	configChangesClientKafka.Default()
}

func randomString(nLetters ...int) string {
	n := 10
	if len(nLetters) > 0 {
		n = nLetters[0]
	}
	//REFERENCE: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func newLogicTest(metaType, protocol string) *logicTest {
	var meta interface {
		internal.Parameterizer
		internal.Configurer
		internal.Initializer
		meta.Employee
	}
	var changesClient interface {
		changesclient.Client
		changesclient.Handler
		internal.Initializer
		internal.Configurer
		internal.Parameterizer
	}
	var changesHandler interface {
		changesclient.Handler
		internal.Initializer
		internal.Configurer
		internal.Parameterizer
	}

	logger := internal_logger.New()
	logger.Configure(&internal_logger.Configuration{
		Level:  internal_logger.Trace,
		Prefix: "bludgeon_rest_server_test",
	})
	switch metaType {
	case "memory":
		meta = memory.New()
	case "file":
		meta = file.New()
	case "mysql":
		meta = mysql.New()
	}
	meta.SetUtilities(logger)
	switch protocol {
	case "rest":
		c := changesclientrest.New()
		changesClient, changesHandler = c, c
	case "kafka":
		c, k := changesclientrest.New(), changesclientkafka.New()
		changesClient, changesHandler = c, k
	}
	changesClient.SetUtilities(logger)
	changesHandler.SetUtilities(logger)
	logic := logic.New()
	logic.SetUtilities(logger)
	logic.SetParameters(meta, changesClient)
	return &logicTest{
		changesClient:  changesClient,
		changesHandler: changesHandler,
		meta:           meta,
		Logic:          logic,
	}
}

func (l *logicTest) initialize(t *testing.T, metaType, protocol string) {
	switch metaType {
	case "file":
		err := l.meta.Configure(configMetaFile)
		assert.Nil(t, err)
	case "mysql":
		err := l.meta.Configure(configMetaMysql)
		assert.Nil(t, err)
	}
	err := l.meta.Initialize()
	assert.Nil(t, err)
	switch protocol {
	case "rest":
		err := l.changesClient.Configure(configChangesClientRest)
		assert.Nil(t, err)
	case "kafka":
		err := l.changesClient.Configure(configChangesClientRest)
		assert.Nil(t, err)
		err = l.changesHandler.Configure(configChangesClientKafka, configKafkaClient)
		assert.Nil(t, err)
	}
	err = l.changesClient.Initialize()
	assert.Nil(t, err)
	err = l.changesHandler.Initialize()
	assert.Nil(t, err)
}

func (l *logicTest) shutdown(t *testing.T) {
	l.meta.Shutdown()
	l.changesClient.Shutdown()
}

func (l *logicTest) TestChanges(t *testing.T) {
	ctx := context.TODO()

	//create employee
	firstName, lastName := randomString(), randomString()
	emailAddress := randomString() + "@name.company"
	employeeCreated, err := l.EmployeeCreate(ctx, data.EmployeePartial{
		FirstName:    &firstName,
		LastName:     &lastName,
		EmailAddress: &emailAddress,
	})
	assert.Nil(t, err)
	assert.NotNil(t, employeeCreated)
	employeeId := employeeCreated.ID
	defer func() {
		l.EmployeeDelete(ctx, employeeId)
		changesRead, _ := l.changesClient.ChangesRead(ctx, changesdata.ChangeSearch{
			DataIds: []string{employeeId},
		})
		for _, change := range changesRead {
			l.changesClient.ChangeDelete(ctx, change.Id)
		}
	}()

	//validate create change
	changesRead, err := l.changesClient.ChangesRead(ctx, changesdata.ChangeSearch{
		DataIds:      []string{employeeId},
		Types:        []string{logic.ChangeTypeEmployee},
		ServiceNames: []string{logic.ServiceName},
		Actions:      []string{logic.ChangeActionCreate},
	})
	assert.Nil(t, err)
	assert.NotNil(t, changesRead)
	assert.Len(t, changesRead, 1)
	assert.Equal(t, employeeCreated.ID, changesRead[0].DataId)
	assert.Equal(t, employeeCreated.Version, changesRead[0].DataVersion)
	assert.Equal(t, logic.ChangeActionCreate, changesRead[0].DataAction)
	assert.Equal(t, logic.ServiceName, changesRead[0].DataServiceName)
	assert.Equal(t, logic.ChangeTypeEmployee, changesRead[0].DataType)

	//update employee
	firstName, lastName = randomString(), randomString()
	employeeUpdated, err := l.EmployeeUpdate(ctx, employeeId, data.EmployeePartial{
		FirstName: &firstName,
		LastName:  &lastName,
	})
	assert.Nil(t, err)
	assert.NotNil(t, employeeUpdated)

	//validate update change
	changesRead, err = l.changesClient.ChangesRead(ctx, changesdata.ChangeSearch{
		DataIds:      []string{employeeId},
		Types:        []string{logic.ChangeTypeEmployee},
		ServiceNames: []string{logic.ServiceName},
		Actions:      []string{logic.ChangeActionUpdate},
	})
	assert.Nil(t, err)
	assert.NotNil(t, changesRead)
	assert.Len(t, changesRead, 1)
	assert.Equal(t, employeeUpdated.ID, changesRead[0].DataId)
	assert.Equal(t, employeeUpdated.Version, changesRead[0].DataVersion)
	assert.Equal(t, logic.ChangeActionUpdate, changesRead[0].DataAction)
	assert.Equal(t, logic.ServiceName, changesRead[0].DataServiceName)
	assert.Equal(t, logic.ChangeTypeEmployee, changesRead[0].DataType)

	// delete employee
	err = l.EmployeeDelete(ctx, employeeId)
	assert.Nil(t, err)

	//validate delete change
	changesRead, err = l.changesClient.ChangesRead(ctx, changesdata.ChangeSearch{
		DataIds:      []string{employeeId},
		Types:        []string{logic.ChangeTypeEmployee},
		ServiceNames: []string{logic.ServiceName},
		Actions:      []string{logic.ChangeActionDelete},
	})
	assert.Nil(t, err)
	assert.NotNil(t, changesRead)
	assert.Len(t, changesRead, 1)
	assert.Equal(t, employeeUpdated.ID, changesRead[0].DataId)
	assert.Equal(t, logic.ChangeActionDelete, changesRead[0].DataAction)
	assert.Equal(t, logic.ServiceName, changesRead[0].DataServiceName)
	assert.Equal(t, logic.ChangeTypeEmployee, changesRead[0].DataType)
}

func testLogic(t *testing.T, metaType, protocol string) {
	l := newLogicTest(metaType, protocol)

	l.initialize(t, metaType, protocol)
	defer l.shutdown(t)

	t.Run("Changes", l.TestChanges)
}

func TestLogicMemoryRest(t *testing.T) {
	testLogic(t, "memory", "rest")
}

func TestLogicFileRest(t *testing.T) {
	testLogic(t, "file", "rest")
}

func TestLogicMysqlRest(t *testing.T) {
	testLogic(t, "mysql", "rest")
}

func TestLogicMemoryKafka(t *testing.T) {
	testLogic(t, "memory", "kafka")
}

func TestLogicFileKafka(t *testing.T) {
	testLogic(t, "file", "kafka")
}

func TestLogicMysqlKafka(t *testing.T) {
	testLogic(t, "mysql", "kafka")
}
