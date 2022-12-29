package service_test

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
	pb "github.com/antonio-alexander/go-bludgeon/timers/data/pb"
	logic "github.com/antonio-alexander/go-bludgeon/timers/logic"
	meta "github.com/antonio-alexander/go-bludgeon/timers/meta"
	file "github.com/antonio-alexander/go-bludgeon/timers/meta/file"
	memory "github.com/antonio-alexander/go-bludgeon/timers/meta/memory"
	mysql "github.com/antonio-alexander/go-bludgeon/timers/meta/mysql"
	service "github.com/antonio-alexander/go-bludgeon/timers/service/grpc"

	changesclient "github.com/antonio-alexander/go-bludgeon/changes/client"
	changesclientkafka "github.com/antonio-alexander/go-bludgeon/changes/client/kafka"
	changesclientrest "github.com/antonio-alexander/go-bludgeon/changes/client/rest"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	internal_server "github.com/antonio-alexander/go-bludgeon/internal/grpc/server"
	internal_kafka "github.com/antonio-alexander/go-bludgeon/internal/kafka"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_meta "github.com/antonio-alexander/go-bludgeon/internal/meta"
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
	configLogger             = new(internal_logger.Configuration)
	configServer             = new(internal_server.Configuration)
	configLogic              = new(logic.Configuration)
	configChangesClientRest  = new(changesclientrest.Configuration)
	configChangesClientKafka = new(changesclientkafka.Configuration)
	configKafkaClient        = new(internal_kafka.Configuration)
)

func randomString(n int) string {
	//REFERENCE: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type grpcServiceTest struct {
	server interface {
		internal.Initializer
		internal.Configurer
	}
	meta interface {
		internal.Initializer
		internal.Configurer
	}
	changesClient interface {
		internal.Initializer
		internal.Configurer
	}
	changesHandler interface {
		internal.Initializer
		internal.Configurer
	}
	logic interface {
		internal.Initializer
		internal.Configurer
	}
	grpcConn *grpc.ClientConn
	pb.TimersClient
	pb.TimeSlicesClient
}

func init() {
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		if s := strings.Split(e, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	configLogger.Default()
	configLogger.FromEnv(envs)
	configMetaFile.Default()
	configMetaFile.FromEnv(envs)
	configMetaFile.File = path.Join("../../tmp", filename)
	os.Remove(configMetaFile.File)
	configMetaMysql.Default()
	configMetaMysql.FromEnv(envs)
	configServer.Default()
	configServer.FromEnv(envs)
	configKafkaClient.Default()
	configKafkaClient.FromEnv(envs)
	configChangesClientKafka.Default()
	configChangesClientRest.Default()
	configChangesClientRest.FromEnv(envs)
	configLogic.Default()
	configLogic.FromEnv(envs)
	configServer.Address = "localhost"
	configServer.Port = "7999"
}

func newGrpcServiceTest(metaType internal_meta.Type, protocol string) *grpcServiceTest {
	var timerMeta interface {
		meta.Timer
		meta.TimeSlice
		internal.Initializer
		internal.Configurer
		internal.Parameterizer
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
		Prefix: "bludgeon_grpc_server_test",
		Level:  internal_logger.Trace,
	})
	switch metaType {
	default:
		timerMeta = memory.New()
	case internal_meta.TypeMySQL:
		timerMeta = mysql.New()
	case internal_meta.TypeFile:
		timerMeta = file.New()
	}
	timerMeta.SetUtilities(logger)
	switch protocol {
	default: //rest
		c := changesclientrest.New()
		changesClient, changesHandler = c, c
	case "kafka":
		c, h := changesclientrest.New(), changesclientkafka.New()
		changesClient, changesHandler = c, h
	}
	changesClient.SetUtilities(logger)
	changesHandler.SetUtilities(logger)
	timerLogic := logic.New()
	timerLogic.SetParameters(timerMeta, changesClient, changesHandler)
	timerLogic.SetUtilities(logger)
	timerService := service.New()
	timerService.SetUtilities(logger)
	timerService.SetParameters(timerLogic)
	server := internal_server.New()
	server.SetUtilities(logger)
	server.SetParameters(timerService)
	return &grpcServiceTest{
		server:           server,
		meta:             timerMeta,
		changesClient:    changesClient,
		changesHandler:   changesHandler,
		logic:            timerLogic,
		grpcConn:         &grpc.ClientConn{},
		TimersClient:     nil,
		TimeSlicesClient: nil,
	}
}

func (r *grpcServiceTest) Initialize(t *testing.T, metaType internal_meta.Type, protocol string) {
	switch metaType {
	case internal_meta.TypeMySQL:
		err := r.meta.Configure(configMetaMysql)
		assert.Nil(t, err)
	case internal_meta.TypeFile:
		err := r.meta.Configure(configMetaFile)
		assert.Nil(t, err)
	}
	switch protocol {
	default: //rest
		err := r.changesClient.Configure(configChangesClientRest)
		assert.Nil(t, err)
		err = r.changesClient.Initialize()
		assert.Nil(t, err)
	case "kafka":
		err := r.changesClient.Configure(configChangesClientRest)
		assert.Nil(t, err)
		err = r.changesClient.Initialize()
		assert.Nil(t, err)
		err = r.changesHandler.Configure(configChangesClientKafka, configKafkaClient)
		assert.Nil(t, err)
		err = r.changesHandler.Initialize()
		assert.Nil(t, err)
	}
	err := r.meta.Initialize()
	assert.Nil(t, err)
	err = r.logic.Configure(configLogic)
	assert.Nil(t, err)
	err = r.logic.Initialize()
	assert.Nil(t, err)
	err = r.server.Configure(configServer)
	assert.Nil(t, err)
	err = r.server.Configure(configServer)
	assert.Nil(t, err)
	err = r.server.Initialize()
	assert.Nil(t, err)
	r.grpcConn, err = grpc.Dial(fmt.Sprintf("%s:%s", configServer.Address, configServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err)
	r.TimersClient = pb.NewTimersClient(r.grpcConn)
	r.TimeSlicesClient = pb.NewTimeSlicesClient(r.grpcConn)
}

func (r *grpcServiceTest) Shutdown(t *testing.T) {
	r.server.Shutdown()
	r.logic.Shutdown()
	r.meta.Shutdown()
	r.grpcConn.Close()
	r.changesClient.Shutdown()
	r.changesHandler.Shutdown()
}

func (r *grpcServiceTest) TestTimerOperations(t *testing.T) {
	ctx := context.TODO()

	//create Timer
	timerCreated, err := r.TimerCreate(ctx, &pb.TimerCreateRequest{
		TimerPartial: pb.FromTimerPartial(&data.TimerPartial{
			//
		}),
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerCreated)
	timerId := timerCreated.GetTimer().Id

	//read created Timer
	timerRead, err := r.TimerRead(ctx, &pb.TimerReadRequest{
		Id: timerId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerRead)
	assert.Equal(t, timerCreated.GetTimer(), timerRead.GetTimer())

	//read all timers
	timersRead, err := r.TimersRead(ctx, &pb.TimersReadRequest{
		TimerSearch: &pb.TimerSearch{
			Ids: []string{timerId},
		},
	})
	assert.Nil(t, err)
	assert.Len(t, timersRead.GetTimers(), 1)
	assert.Contains(t, timersRead.GetTimers(), timerCreated.GetTimer())

	//start timer
	timerStarted, err := r.TimerStart(ctx, &pb.TimerStartRequest{
		Id: timerId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerStarted.GetTimer())
	assert.NotZero(t, timerStarted.GetTimer().GetStart())

	time.Sleep(time.Second)

	//stop timer
	timerStopped, err := r.TimerStop(ctx, &pb.TimerStopRequest{
		Id: timerId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerStopped.GetTimer())
	assert.GreaterOrEqual(t, timerStopped.GetTimer().GetElapsedTime(), time.Second)

	//read timer
	timerRead, err = r.TimerRead(ctx, &pb.TimerReadRequest{
		Id: timerId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerRead)
	assert.Equal(t, timerStopped.GetTimer(), timerRead.GetTimer())
	assert.Equal(t, timerStopped.GetTimer().GetElapsedTime(), timerRead.GetTimer().GetElapsedTime())

	//submit timer
	tNow := time.Now()
	timerSubmitted, err := r.TimerSubmit(ctx, &pb.TimerSubmitRequest{
		Id: timerId,
		FinishOneof: &pb.TimerSubmitRequest_Finish{
			Finish: tNow.UnixNano(),
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerSubmitted.GetTimer())

	//update timer comment
	comment := randomString(25)
	timerUpdated, err := r.TimerUpdate(ctx, &pb.TimerUpdateRequest{
		Id: timerId,
		TimerPartial: pb.FromTimerPartial(&data.TimerPartial{
			Comment: &comment,
		}),
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerUpdated)
	assert.Equal(t, comment, timerUpdated.GetTimer().GetComment())

	//read updated Timer
	timerRead, err = r.TimerRead(ctx, &pb.TimerReadRequest{
		Id: timerId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerRead)
	assert.Equal(t, timerUpdated.GetTimer(), timerRead.GetTimer())

	//archive timer
	archived := true
	timerUpdated, err = r.TimerUpdate(ctx, &pb.TimerUpdateRequest{
		Id: timerId,
		TimerPartial: pb.FromTimerPartial(&data.TimerPartial{
			Archived: &archived,
		}),
	})
	assert.Nil(t, err)
	assert.True(t, timerUpdated.GetTimer().GetArchived())

	//delete Timer
	_, err = r.TimerDelete(ctx, &pb.TimerDeleteRequest{
		Id: timerId,
	})
	assert.Nil(t, err)

	//delete Timer again
	_, err = r.TimerDelete(ctx, &pb.TimerDeleteRequest{
		Id: timerId,
	})
	assert.NotNil(t, err)

	//attempt to read deleted Timer
	timerRead, err = r.TimerRead(ctx, &pb.TimerReadRequest{
		Id: timerId,
	})
	assert.NotNil(t, err)
	assert.Nil(t, timerRead.GetTimer())
}

func testTimersGrpcService(t *testing.T, metaType internal_meta.Type, protocol string) {
	r := newGrpcServiceTest(metaType, protocol)
	r.Initialize(t, metaType, protocol)
	defer r.Shutdown(t)

	t.Run("Test Timer Operations", r.TestTimerOperations)
}

func TestTimersGrpcServiceMemory(t *testing.T) {
	testTimersGrpcService(t, internal_meta.TypeMemory, "rest")
}

func TestTimersGrpcServiceFile(t *testing.T) {
	testTimersGrpcService(t, internal_meta.TypeFile, "rest")
}

func TestTimersGrpcServiceMysql(t *testing.T) {
	testTimersGrpcService(t, internal_meta.TypeMySQL, "rest")
}
