package service_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/timers/data"
	"github.com/antonio-alexander/go-bludgeon/timers/data/pb"
	"github.com/antonio-alexander/go-bludgeon/timers/logic"
	"github.com/antonio-alexander/go-bludgeon/timers/meta"
	"github.com/antonio-alexander/go-bludgeon/timers/meta/memory"
	service "github.com/antonio-alexander/go-bludgeon/timers/service/grpc"

	internal_server "github.com/antonio-alexander/go-bludgeon/internal/grpc/server"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_meta "github.com/antonio-alexander/go-bludgeon/internal/meta"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	address string            = "localhost"
	port    string            = "8081"
	options []grpc.DialOption = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
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
	server  internal_server.Owner
	service service.Owner
	meta    interface {
		meta.Serializer
		meta.Timer
		internal_meta.Owner
	}
	logic interface {
		logic.Logic
	}
	conn   *grpc.ClientConn
	client pb.TimersClient
}

func init() {
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		if s := strings.Split(e, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	if _, ok := envs["BLUDGEON_GRPC_ADDRESS"]; ok {
		address = envs["BLUDGEON_GRPC_ADDRESS"]
	}
	if _, ok := envs["BLUDGEON_GRPC_PORT"]; ok {
		port = envs["BLUDGEON_GRPC_PORT"]
	}
}

func new() *grpcServiceTest {
	logger := internal_logger.New("bludgeon_grpc_server_test")
	conn, _ := grpc.Dial(fmt.Sprintf("%s:%s", address, port), options...)
	meta := memory.New()
	logic := logic.New(logger, meta)
	client := pb.NewTimersClient(conn)
	server := internal_server.New(logger)
	service := service.New(logger, logic, server)
	return &grpcServiceTest{
		server:  server,
		meta:    meta,
		logic:   logic,
		client:  client,
		conn:    conn,
		service: service,
	}
}

func (r *grpcServiceTest) initialize(t *testing.T) {
	err := r.server.Initialize(&internal_server.Configuration{
		Address: address,
		Port:    port,
		Options: []grpc.ServerOption{},
	}, r.service.Register)
	assert.Nil(t, err)
}

func (r *grpcServiceTest) shutdown(t *testing.T) {
	r.meta.Shutdown()
	r.server.Shutdown()
}

func (r *grpcServiceTest) TestTimerOperations(t *testing.T) {
	employeeId := randomString(25)
	ctx := context.TODO()

	//create Timer
	timerCreated, err := r.client.TimerCreate(ctx, &pb.TimerCreateRequest{
		TimerPartial: pb.FromTimerPartial(&data.TimerPartial{
			EmployeeID: &employeeId,
		}),
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerCreated)
	timerId := timerCreated.GetTimer().Id

	//read created Timer
	timerRead, err := r.client.TimerRead(ctx, &pb.TimerReadRequest{
		Id: timerId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerRead)
	assert.Equal(t, timerCreated.GetTimer(), timerRead.GetTimer())

	//read all timers
	timersRead, err := r.client.TimersRead(ctx, &pb.TimersReadRequest{
		TimerSearch: &pb.TimerSearch{},
	})
	assert.Nil(t, err)
	assert.Len(t, timersRead.GetTimers(), 1)
	assert.Contains(t, timersRead.GetTimers(), timerCreated.GetTimer())

	//start timer
	timerStarted, err := r.client.TimerStart(ctx, &pb.TimerStartRequest{
		Id: timerId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerStarted.GetTimer())
	assert.NotZero(t, timerStarted.GetTimer().GetStart())

	time.Sleep(time.Second)

	//stop timer
	timerStopped, err := r.client.TimerStop(ctx, &pb.TimerStopRequest{
		Id: timerId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerStopped.GetTimer())
	assert.GreaterOrEqual(t, timerStopped.GetTimer().GetElapsedTime(), int64(time.Second))

	//read timer
	timerRead, err = r.client.TimerRead(ctx, &pb.TimerReadRequest{
		Id: timerId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerRead)
	assert.Equal(t, timerStopped.GetTimer(), timerRead.GetTimer())
	assert.Equal(t, timerStopped.GetTimer().GetElapsedTime(), timerRead.GetTimer().GetElapsedTime())

	//submit timer
	tNow := time.Now()
	timerSubmitted, err := r.client.TimerSubmit(ctx, &pb.TimerSubmitRequest{
		Id: timerId,
		FinishOneof: &pb.TimerSubmitRequest_Finish{
			Finish: tNow.UnixNano(),
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerSubmitted.GetTimer())
	assert.Equal(t, tNow.UnixNano(), timerSubmitted.GetTimer().GetFinish())
	assert.True(t, timerSubmitted.GetTimer().GetCompleted())

	//update timer comment
	comment := randomString(25)
	timerUpdated, err := r.client.TimerUpdateComment(ctx, &pb.TimerUpdateCommentRequest{
		Id:      timerId,
		Comment: comment,
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerUpdated)
	assert.Equal(t, comment, timerUpdated.GetTimer().GetComment())

	//read updated Timer
	timerRead, err = r.client.TimerRead(ctx, &pb.TimerReadRequest{
		Id: timerId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, timerRead)
	assert.Equal(t, timerUpdated.GetTimer(), timerRead.GetTimer())

	//TODO: archive timer
	timerArchived, err := r.client.TimerArchive(ctx, &pb.TimerArchiveRequest{
		Id:      timerId,
		Archive: true,
	})
	assert.Nil(t, err)
	assert.True(t, timerArchived.GetTimer().GetArchived())

	//delete Timer
	_, err = r.client.TimerDelete(ctx, &pb.TimerDeleteRequest{
		Id: timerId,
	})
	assert.Nil(t, err)

	//delete Timer again
	_, err = r.client.TimerDelete(ctx, &pb.TimerDeleteRequest{
		Id: timerId,
	})
	assert.NotNil(t, err)

	//attempt to read deleted Timer
	timerRead, err = r.client.TimerRead(ctx, &pb.TimerReadRequest{
		Id: timerId,
	})
	assert.NotNil(t, err)
	assert.Nil(t, timerRead.GetTimer())
}

func TestTimersGrpcService(t *testing.T) {
	r := new()
	r.initialize(t)
	t.Run("Test Timer Operations", r.TestTimerOperations)
	r.shutdown(t)
}
