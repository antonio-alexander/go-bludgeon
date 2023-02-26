package logic_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	data "github.com/antonio-alexander/go-bludgeon/timers/data"
	logic "github.com/antonio-alexander/go-bludgeon/timers/logic"
	meta "github.com/antonio-alexander/go-bludgeon/timers/meta"
	file "github.com/antonio-alexander/go-bludgeon/timers/meta/file"
	memory "github.com/antonio-alexander/go-bludgeon/timers/meta/memory"
	mysql "github.com/antonio-alexander/go-bludgeon/timers/meta/mysql"

	changesclient "github.com/antonio-alexander/go-bludgeon/changes/client"
	changesclientkafka "github.com/antonio-alexander/go-bludgeon/changes/client/kafka"
	changesclientrest "github.com/antonio-alexander/go-bludgeon/changes/client/rest"
	changesdata "github.com/antonio-alexander/go-bludgeon/changes/data"
	employeesclient "github.com/antonio-alexander/go-bludgeon/employees/client"
	employeesclientgrpc "github.com/antonio-alexander/go-bludgeon/employees/client/grpc"
	employeesclientrest "github.com/antonio-alexander/go-bludgeon/employees/client/rest"
	employeesdata "github.com/antonio-alexander/go-bludgeon/employees/data"
	_ "github.com/antonio-alexander/go-bludgeon/employees/data/pb"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	internal_kafka "github.com/antonio-alexander/go-bludgeon/internal/kafka"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const filename string = "bludgeon_logic.json"

var (
	configMetaMysql          = new(internal_mysql.Configuration)
	configMetaFile           = new(internal_file.Configuration)
	configChangesClientRest  = new(changesclientrest.Configuration)
	configChangesClientKafka = new(changesclientkafka.Configuration)
	configEmployeeClientRest = new(employeesclientrest.Configuration)
	configEmployeeClientGrpc = new(employeesclientgrpc.Configuration)
	configKafkaClient        = new(internal_kafka.Configuration)
	configLogic              = new(logic.Configuration)
	letterRunes              = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type logicTest struct {
	meta interface {
		internal.Initializer
		internal.Configurer
	}
	changesClient interface {
		internal.Initializer
		internal.Configurer
		changesclient.Client
	}
	changesHandler interface {
		internal.Initializer
		internal.Configurer
	}
	logic interface {
		internal.Initializer
		internal.Configurer
	}
	employeesClient interface {
		internal.Configurer
		internal.Parameterizer
		internal.Initializer
		employeesclient.Client
	}
	logic.Logic
}

func init() {
	rand.Seed(time.Now().UnixNano())
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
	configChangesClientRest.Default()
	configChangesClientRest.FromEnv(envs)
	configChangesClientRest.Rest.Timeout = 30 * time.Second
	configKafkaClient.Default()
	configChangesClientKafka.Default()
	configLogic.Default()
	configLogic.FromEnv(envs)
	configLogic.ChangeRateRead = time.Second
	configLogic.ChangesTimeout = 30 * time.Second
	//KIM: if we use the default it could conflict with the
	// the timers service/container
	configLogic.ChangesRegistrationId = randomString()
	configEmployeeClientRest.Default()
	configEmployeeClientRest.FromEnv(envs)
	configEmployeeClientGrpc.Default()
	configEmployeeClientGrpc.FromEnv(envs)
	configMetaFile.File = path.Join("../../tmp", filename)
	os.Remove(configMetaFile.File)
	configKafkaClient.Brokers = []string{"localhost:9092"}
	configEmployeeClientGrpc.Options = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
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
		meta.Timer
		meta.TimeSlice
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
	var employeesClient interface {
		internal.Configurer
		internal.Parameterizer
		internal.Initializer
		employeesclient.Client
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
	default: //rest
		c := changesclientrest.New()
		changesClient, changesHandler = c, c
	case "kafka":
		c, h := changesclientrest.New(), changesclientkafka.New()
		changesClient, changesHandler = c, h
	}
	switch protocol {
	default: //rest
		employeesClient = employeesclientrest.New()
	case "grpc":
		employeesClient = employeesclientgrpc.New()
	}
	employeesClient.SetUtilities(logger)
	changesClient.SetUtilities(logger)
	changesHandler.SetUtilities(logger)
	logic := logic.New()
	logic.SetUtilities(logger)
	logic.SetParameters(meta, changesClient, changesHandler)
	return &logicTest{
		meta:            meta,
		logic:           logic,
		changesClient:   changesClient,
		changesHandler:  changesHandler,
		employeesClient: employeesClient,
		Logic:           logic,
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
	default: //rest
		err := l.changesClient.Configure(configChangesClientRest)
		assert.Nil(t, err)
	case "kafka":
		err := l.changesClient.Configure(configChangesClientRest)
		assert.Nil(t, err)
		err = l.changesHandler.Configure(configChangesClientKafka, configKafkaClient)
		assert.Nil(t, err)
	}
	switch protocol {
	default: //rest
		err = l.employeesClient.Configure(configEmployeeClientRest)
		assert.Nil(t, err)
	case "grpc":
		err = l.employeesClient.Configure(configEmployeeClientGrpc)
		assert.Nil(t, err)
	}
	err = l.employeesClient.Initialize()
	assert.Nil(t, err)
	err = l.changesClient.Initialize()
	assert.Nil(t, err)
	err = l.changesHandler.Initialize()
	assert.Nil(t, err)
	err = l.logic.Configure(configLogic)
	assert.Nil(t, err)
	err = l.logic.Initialize()
	assert.Nil(t, err)
	//block until logic connected
	tCheck := time.NewTicker(time.Second)
	defer tCheck.Stop()
	tStop := time.After(30 * time.Second)
	for stop := false; !stop; {
		select {
		case <-tCheck.C:
			if l.IsConnected() {
				stop = true
			}
		case <-tStop:
			assert.Fail(t, "unable to confirm logic connected")
			stop = true
		}
	}
}

func (l *logicTest) shutdown(t *testing.T) {
	l.logic.Shutdown()
	l.meta.Shutdown()
	l.changesClient.RegistrationDelete(context.TODO(), configLogic.ChangesRegistrationId)
	l.changesClient.Shutdown()
}

func (l *logicTest) assertTimerChange(t *testing.T, ctx context.Context, timer *data.Timer, action string) func() bool {
	return func() bool {
		checkChangesFx := func() (string, bool) {
			changesRead, err := l.changesClient.ChangesRead(ctx, changesdata.ChangeSearch{
				DataIds:      []string{timer.ID},
				Types:        []string{data.ChangeTypeTimer},
				ServiceNames: []string{data.ServiceName},
				Actions:      []string{action},
			})
			if err != nil {
				return err.Error(), false
			}
			if len(changesRead) != 1 {
				return fmt.Sprintf("Number of changes equal to: %d", len(changesRead)), false
			}
			if changesRead[0].DataId != timer.ID {
				return fmt.Sprintf("DataId != %s", timer.ID), false
			}
			if changesRead[0].DataAction != data.ChangeActionDelete &&
				changesRead[0].DataVersion != timer.Version {
				return fmt.Sprintf("Version not equal to %d", timer.Version), false
			}
			return "", true
		}
		tCheck := time.NewTicker(time.Second)
		defer tCheck.Stop()
		tStop := time.After(10 * time.Second)
		for {
			select {
			case <-tStop:
				reason, success := checkChangesFx()
				if !success {
					t.Logf("failure because of: %s", reason)
				}
				return false
			case <-tCheck.C:
				if _, success := checkChangesFx(); success {
					return success
				}
			}
		}
	}
}

func (l *logicTest) TestTimerChanges(t *testing.T) {
	ctx := context.TODO()

	//create timer
	comment := randomString(25)
	timerCreated, err := l.TimerCreate(ctx, data.TimerPartial{
		Comment: &comment,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, timerCreated.ID)
	assert.Equal(t, timerCreated.Comment, comment)
	timerId := timerCreated.ID
	defer func() {
		l.TimerDelete(ctx, timerId)
		changesRead, _ := l.changesClient.ChangesRead(ctx, changesdata.ChangeSearch{
			DataIds: []string{timerId},
		})
		for _, change := range changesRead {
			l.changesClient.ChangeDelete(ctx, change.Id)
		}
	}()

	//wait for change to be sent (asynchronous)
	time.Sleep(time.Second)

	//validate create change
	assert.Condition(t, l.assertTimerChange(t, ctx, timerCreated, data.ChangeActionCreate))

	//update timer
	comment = randomString()
	timerUpdated, err := l.TimerUpdate(ctx, timerId, data.TimerPartial{
		Comment: &comment,
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerUpdated)

	//wait for change to be sent (asynchronous)
	time.Sleep(time.Second)

	//validate update change
	assert.Condition(t, l.assertTimerChange(t, ctx, timerUpdated, data.ChangeActionUpdate))

	// start timer
	timerStarted, err := l.TimerStart(ctx, timerId)
	assert.Nil(t, err)
	assert.NotNil(t, timerStarted)

	//wait for change to be sent (asynchronous)
	time.Sleep(time.Second)

	//validate  change
	assert.Condition(t, l.assertTimerChange(t, ctx, timerStarted, data.ChangeActionStart))

	//sleep to give timer a duration/elapsed time
	time.Sleep(time.Second)

	// stop timer
	timerStopped, err := l.TimerStop(ctx, timerId)
	assert.Nil(t, err)
	assert.NotNil(t, timerStopped)

	//wait for change to be sent (asynchronous)
	time.Sleep(time.Second)

	//validate  change
	assert.Condition(t, l.assertTimerChange(t, ctx, timerStopped, data.ChangeActionStop))

	// submit timer
	tNow := time.Now()
	timerSubmitted, err := l.TimerSubmit(ctx, timerId, tNow.UnixNano())
	assert.Nil(t, err)
	assert.NotNil(t, timerSubmitted)

	//wait for change to be sent (asynchronous)
	time.Sleep(time.Second)

	//validate  change
	assert.Condition(t, l.assertTimerChange(t, ctx, timerSubmitted, data.ChangeActionSubmit))

	//delete timer
	err = l.TimerDelete(ctx, timerId)
	assert.Nil(t, err)

	//wait for change to be sent (asynchronous)
	time.Sleep(time.Second)

	//validate  change
	assert.Condition(t, l.assertTimerChange(t, ctx, timerSubmitted, data.ChangeActionDelete))
}

func (l *logicTest) TestEmployeeChanges(t *testing.T) {
	ctx := context.TODO()

	//create employee
	firstName, lastName := randomString(), randomString()
	emailAddress := randomString() + "@foobar.duck"
	employeeCreated, err := l.employeesClient.EmployeeCreate(ctx, employeesdata.EmployeePartial{
		FirstName:    &firstName,
		LastName:     &lastName,
		EmailAddress: &emailAddress,
	})
	assert.Nil(t, err)
	assert.NotNil(t, employeeCreated)
	employeeId := employeeCreated.ID
	defer func() {
		l.employeesClient.EmployeeDelete(ctx, employeeId)
	}()

	//create timer
	comment := randomString(25)
	timerCreated, err := l.TimerCreate(ctx, data.TimerPartial{
		Comment:    &comment,
		EmployeeID: &employeeId,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, timerCreated.ID)
	assert.Equal(t, timerCreated.Comment, comment)
	assert.Equal(t, employeeId, timerCreated.EmployeeID)
	timerId := timerCreated.ID
	defer func() {
		l.TimerDelete(ctx, timerId)
		changesRead, _ := l.changesClient.ChangesRead(ctx, changesdata.ChangeSearch{
			DataIds: []string{timerId},
		})
		for _, change := range changesRead {
			l.changesClient.ChangeDelete(ctx, change.Id)
		}
	}()

	// delete employee
	err = l.employeesClient.EmployeeDelete(ctx, employeeId)
	assert.Nil(t, err)

	//wait for change to be read by timers service
	time.Sleep(time.Second)

	//validate that timer deleted
	tCheck := time.NewTicker(time.Second)
	defer tCheck.Stop()
	tStop := time.After(10 * time.Second)
	for stop := false; !stop; {
		select {
		case <-tStop:
			assert.Fail(t, "unable to confirm timer deleted")
			stop = true
		case <-tCheck.C:
			timerRead, err := l.TimerRead(ctx, timerId)
			switch {
			case err != nil:
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, meta.ErrTimerNotFound)
				assert.Nil(t, timerRead)
				stop = true
			default:
				assert.NotNil(t, timerRead)
			}
		}
	}
}

func testLogic(t *testing.T, metaType, protocol string) {
	l := newLogicTest(metaType, protocol)

	l.initialize(t, metaType, protocol)
	defer l.shutdown(t)

	//wait for registration to occur
	time.Sleep(5 * time.Second)

	t.Run("Employee Changes", l.TestEmployeeChanges)
	t.Run("Timer Changes", l.TestTimerChanges)

	//sleep to ensure separation between tests
	time.Sleep(5 * time.Second)
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

func TestLogicMemoryGrpc(t *testing.T) {
	testLogic(t, "memory", "grpc")
}

func TestLogicFileGrpc(t *testing.T) {
	testLogic(t, "file", "grpc")
}

func TestLogicMysqlGrpc(t *testing.T) {
	testLogic(t, "mysql", "grpc")
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
