package service_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"testing"

	data "github.com/antonio-alexander/go-bludgeon/employees/data"
	pb "github.com/antonio-alexander/go-bludgeon/employees/data/pb"
	logic "github.com/antonio-alexander/go-bludgeon/employees/logic"
	meta "github.com/antonio-alexander/go-bludgeon/employees/meta"
	file "github.com/antonio-alexander/go-bludgeon/employees/meta/file"
	memory "github.com/antonio-alexander/go-bludgeon/employees/meta/memory"
	mysql "github.com/antonio-alexander/go-bludgeon/employees/meta/mysql"
	service "github.com/antonio-alexander/go-bludgeon/employees/service/grpc"

	changesclient "github.com/antonio-alexander/go-bludgeon/changes/client"
	changesclientrest "github.com/antonio-alexander/go-bludgeon/changes/client/rest"

	internal "github.com/antonio-alexander/go-bludgeon/internal"
	internal_server "github.com/antonio-alexander/go-bludgeon/internal/grpc/server"
	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_file "github.com/antonio-alexander/go-bludgeon/internal/meta/file"
	internal_mysql "github.com/antonio-alexander/go-bludgeon/internal/meta/mysql"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const filename string = "bludgeon_logic.json"

var (
	configMetaMysql         = new(internal_mysql.Configuration)
	configMetaFile          = new(internal_file.Configuration)
	configLogger            = new(internal_logger.Configuration)
	configServer            = new(internal_server.Configuration)
	configLogic             = new(logic.Configuration)
	configChangesClientRest = new(changesclientrest.Configuration)
	letterRunes             = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
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
	server interface {
		internal.Initializer
		internal.Configurer
	}
	meta interface {
		internal.Initializer
		internal.Configurer
	}
	logic interface {
		logic.Logic
	}
	logger interface {
		internal.Configurer
	}
	changesClient interface {
		internal.Initializer
		internal.Configurer
		changesclient.Client
	}
	grpcConn *grpc.ClientConn
	pb.EmployeesClient
}

func init() {
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		if s := strings.Split(e, "="); len(s) > 1 {
			envs[s[0]] = strings.Join(s[1:], "=")
		}
	}
	// options     []grpc.DialOption =
	configLogger.Default()
	configLogger.FromEnv(envs)
	configMetaFile.Default()
	configMetaFile.FromEnv(envs)
	configMetaFile.File = path.Join("../../tmp", filename)
	os.Remove(configMetaFile.File)
	configMetaMysql.Default()
	configMetaMysql.FromEnv(envs)
	configChangesClientRest.Default()
	configChangesClientRest.FromEnv(envs)
	configLogic.Default()
	configLogic.FromEnv(envs)
	configServer.Default()
	configServer.FromEnv(envs)
	configServer.Address = "localhost"
	configServer.Port = "8082"
}

func newGrpcServiceTest(metaType, protocol string) *grpcServiceTest {
	var meta interface {
		meta.Employee
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

	logger := internal_logger.New()
	logger.Configure(&internal_logger.Configuration{
		Prefix: "bludgeon_grpc_server_test",
		Level:  internal_logger.Trace,
	})
	switch metaType {
	default:
		meta = memory.New()
	case "mysql":
		meta = mysql.New()
	case "file":
		meta = file.New()
	}
	meta.SetUtilities(logger)
	switch protocol {
	case "rest":
		changesClient = changesclientrest.New()
	}
	changesClient.SetUtilities(logger)
	logic := logic.New()
	logic.SetParameters(meta, changesClient)
	logic.SetUtilities(logger)
	logic.Configure(configLogic)
	service := service.New()
	service.SetUtilities(logger)
	service.SetParameters(logic)
	server := internal_server.New()
	server.SetUtilities(logger)
	server.SetParameters(logic, service)
	return &grpcServiceTest{
		server:        server,
		meta:          meta,
		logic:         logic,
		logger:        logger,
		changesClient: changesClient,
	}
}

func (r *grpcServiceTest) initialize(t *testing.T, metaType, protocol string) {
	switch metaType {
	case "file":
		err := r.meta.Configure(configMetaFile)
		assert.Nil(t, err)
	case "mysql":
		err := r.meta.Configure(configMetaMysql)
		assert.Nil(t, err)
	}
	err := r.meta.Initialize()
	assert.Nil(t, err)
	switch protocol {
	case "rest":
		err := r.changesClient.Configure(configChangesClientRest)
		assert.Nil(t, err)
	}
	err = r.changesClient.Initialize()
	assert.Nil(t, err)
	err = r.server.Configure(configServer)
	assert.Nil(t, err)
	err = r.server.Initialize()
	if !assert.Nil(t, err) {
		fmt.Println(err)
	}
	r.grpcConn, err = grpc.Dial(fmt.Sprintf("%s:%s", configServer.Address, configServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err)
	r.EmployeesClient = pb.NewEmployeesClient(r.grpcConn)
}

func (r *grpcServiceTest) shutdown(t *testing.T) {
	r.server.Shutdown()
	r.meta.Shutdown()
	r.changesClient.Shutdown()
	r.grpcConn.Close()
}

func (r *grpcServiceTest) testEmployeeOperations(t *testing.T) {
	firstName, lastName := randomString(25), randomString(25)
	emailAddress := fmt.Sprintf("%s@foobar.duck", randomString(25))
	ctx := context.TODO()

	//create employee
	employeeCreated, err := r.EmployeeCreate(ctx, &pb.EmployeeCreateRequest{
		EmployeePartial: pb.FromEmployeePartial(&data.EmployeePartial{
			FirstName:    &firstName,
			LastName:     &lastName,
			EmailAddress: &emailAddress,
		}),
	})
	assert.Nil(t, err)
	assert.NotNil(t, employeeCreated)
	employeeId := employeeCreated.GetEmployee().Id
	defer func() {
		r.EmployeeDelete(ctx, &pb.EmployeeDeleteRequest{
			Id: employeeId,
		})
	}()

	//read created employee
	employeeRead, err := r.EmployeeRead(ctx, &pb.EmployeeReadRequest{
		Id: employeeId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, employeeRead)
	assert.Equal(t, employeeCreated.GetEmployee(), employeeRead.GetEmployee())

	//read all employees
	employeesRead, err := r.EmployeesRead(ctx, &pb.EmployeesReadRequest{
		EmployeeSearch: &pb.EmployeeSearch{
			Ids: []string{employeeId},
		},
	})
	assert.Nil(t, err)
	assert.Len(t, employeesRead.GetEmployees(), 1)
	assert.Contains(t, employeesRead.GetEmployees(), employeeCreated.GetEmployee())

	//update employee
	updatedFirstName := randomString(25)
	employeeUpdated, err := r.EmployeeUpdate(ctx, &pb.EmployeeUpdateRequest{
		Id: employeeId,
		EmployeePartial: &pb.EmployeePartial{
			FirstNameOneof: &pb.EmployeePartial_FirstName{FirstName: updatedFirstName},
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, employeeUpdated)
	assert.Equal(t, updatedFirstName, employeeUpdated.Employee.GetFirstName())

	//read updated employee
	employeeRead, err = r.EmployeeRead(ctx, &pb.EmployeeReadRequest{
		Id: employeeId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, employeeRead)
	assert.Equal(t, employeeUpdated.GetEmployee(), employeeRead.GetEmployee())

	//delete employee
	_, err = r.EmployeeDelete(ctx, &pb.EmployeeDeleteRequest{
		Id: employeeId,
	})
	assert.Nil(t, err)

	//delete employee again
	_, err = r.EmployeeDelete(ctx, &pb.EmployeeDeleteRequest{
		Id: employeeId,
	})
	assert.NotNil(t, err)

	//attempt to read deleted employee
	employeeRead, err = r.EmployeeRead(ctx, &pb.EmployeeReadRequest{
		Id: employeeId,
	})
	assert.NotNil(t, err)
	assert.Nil(t, employeeRead.GetEmployee())
}

func testEmployeesGrpcService(t *testing.T, metaType string) {
	r := newGrpcServiceTest(metaType, "rest")

	r.initialize(t, metaType, "rest")
	defer r.shutdown(t)

	t.Run("Test Employee Operations", r.testEmployeeOperations)
}

func TestEmployeesGrpcServiceMemory(t *testing.T) {
	testEmployeesGrpcService(t, "memory")
}

func TestEmployeesGrpcServiceFile(t *testing.T) {
	testEmployeesGrpcService(t, "file")
}

func TestEmployeesGrpcServiceMysql(t *testing.T) {
	testEmployeesGrpcService(t, "mysql")
}
