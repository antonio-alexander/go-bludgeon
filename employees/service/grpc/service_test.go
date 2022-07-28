package service_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	pb "github.com/antonio-alexander/go-bludgeon/employees/data/pb"
	logic "github.com/antonio-alexander/go-bludgeon/employees/logic"
	meta "github.com/antonio-alexander/go-bludgeon/employees/meta"
	memory "github.com/antonio-alexander/go-bludgeon/employees/meta/memory"
	service "github.com/antonio-alexander/go-bludgeon/employees/service/grpc"

	internal_server "github.com/antonio-alexander/go-bludgeon/internal/grpc/server"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_meta "github.com/antonio-alexander/go-bludgeon/internal/meta"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	address     string            = "localhost"
	port        string            = "8081"
	options     []grpc.DialOption = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	letterRunes []rune            = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func randomString(n int) string {
	//REFERENCE: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
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
		internal_meta.Owner
		meta.Serializer
		meta.Employee
	}
	logic interface {
		logic.Logic
	}
	conn   *grpc.ClientConn
	client pb.EmployeesClient
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
	client := pb.NewEmployeesClient(conn)
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

func (r *grpcServiceTest) testEmployeeOperations(t *testing.T) {
	firstName, lastName := randomString(25), randomString(25)
	emailAddress := fmt.Sprintf("%s@foobar.duck", randomString(25))
	ctx := context.TODO()

	//create employee
	employeeCreated, err := r.client.EmployeeCreate(ctx, &pb.EmployeeCreateRequest{
		EmployeePartial: pb.FromEmployeePartial(&data.EmployeePartial{
			FirstName:    &firstName,
			LastName:     &lastName,
			EmailAddress: &emailAddress,
		}),
	})
	assert.Nil(t, err)
	assert.NotNil(t, employeeCreated)
	employeeId := employeeCreated.GetEmployee().Id

	//read created employee
	employeeRead, err := r.client.EmployeeRead(ctx, &pb.EmployeeReadRequest{
		Id: employeeId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, employeeRead)
	assert.Equal(t, employeeCreated.GetEmployee(), employeeRead.GetEmployee())

	//read all employees
	employeesRead, err := r.client.EmployeesRead(ctx, &pb.EmployeesReadRequest{
		EmployeeSearch: &pb.EmployeeSearch{},
	})
	assert.Nil(t, err)
	assert.Len(t, employeesRead.GetEmployees(), 1)
	assert.Contains(t, employeesRead.GetEmployees(), employeeCreated.GetEmployee())

	//update employee
	updatedFirstName := randomString(25)
	employeeUpdated, err := r.client.EmployeeUpdate(ctx, &pb.EmployeeUpdateRequest{
		Id: employeeId,
		EmployeePartial: &pb.EmployeePartial{
			FirstNameOneof: &pb.EmployeePartial_FirstName{FirstName: updatedFirstName},
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, employeeUpdated)
	assert.Equal(t, updatedFirstName, employeeUpdated.Employee.GetFirstName())

	//read updated employee
	employeeRead, err = r.client.EmployeeRead(ctx, &pb.EmployeeReadRequest{
		Id: employeeId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, employeeRead)
	assert.Equal(t, employeeUpdated.GetEmployee(), employeeRead.GetEmployee())

	//delete employee
	_, err = r.client.EmployeeDelete(ctx, &pb.EmployeeDeleteRequest{
		Id: employeeId,
	})
	assert.Nil(t, err)

	//delete employee again
	_, err = r.client.EmployeeDelete(ctx, &pb.EmployeeDeleteRequest{
		Id: employeeId,
	})
	assert.NotNil(t, err)

	//attempt to read deleted employee
	employeeRead, err = r.client.EmployeeRead(ctx, &pb.EmployeeReadRequest{
		Id: employeeId,
	})
	assert.NotNil(t, err)
	assert.Nil(t, employeeRead.GetEmployee())
}

func TestEmployeesGrpcService(t *testing.T) {
	r := new()
	r.initialize(t)
	t.Run("Test Employee Operations", r.testEmployeeOperations)
	r.shutdown(t)
}
