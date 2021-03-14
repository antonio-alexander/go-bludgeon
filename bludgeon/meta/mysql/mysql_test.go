package bludgeonmetamysql_test

//--------------------------------------------------------------------------------------------
// database_test.go contains all the tests to verify functionality of the bludgeon-database
// library, it contains all the unit and functions tests specific to the database
//--------------------------------------------------------------------------------------------

import (
	"testing"
	"time"

	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/mysql"
	"github.com/stretchr/testify/assert"
)

//--------------------------------------------------------------------------------------------------
//
//
// Normal Use Cases:
//
// Edge Cases:
//
//--------------------------------------------------------------------------------------------------

const (
	TestDatabaseName string = "bludgeon"
	rootUsername     string = "root"
	bludgeonUsername string = "bludgeon"
)

const (
	testCaseMap string = "Test case: %s"
)

var (
	hostname string = mysql.DefaultHostname
	port     string = mysql.DefaultPort
	username string = mysql.DefaultUsername
	password string = mysql.DefaultPassword
)

func init() {
	//TODO: setup variables from environment?
}

//--------------------------------------------------------------------------------------------------
// UNIT TESTS
// Purpose: Unit Tests can only check the input and output of exported functions. For cases, inputs
// can be prefixed with an 'i' and outputs with an 'o. Use a map that uses a string and an anonymous
// struct. The string is the case description and the struct is a collection of inputs and outputs
//
// Function Prefix: TestUnit
//--------------------------------------------------------------------------------------------------

//--------------------------------------------------------------------------------------------------
// FUNCTION TESTS
//
// Purpose: Function Tests check the use of multiple package functions that do not rely on an
// external source
// Function Prefix: TestFunc
//
// Progression:
// 1. Level 1
// 		a. Level 2
// 			(1) Level 3
//--------------------------------------------------------------------------------------------------

//--------------------------------------------------------------------------------------------------
// INTEGRATION TESTS
//
// Purpose: Integration tests check the use of multiple package functions that rely on one or more
// external source
// Function Prefix: TestInt
//
// Progression:
// 1. Level 1
// 		a. Level 2
// 			(1) Level 3
//--------------------------------------------------------------------------------------------------

func TestIntConnect(t *testing.T) {
	//Test: this unit test is meant to test whether or not the connect function works and to validate
	// certain use cases for that connect function
	//Notes:
	//Verification:

	cases := map[string]struct {
		iConfig mysql.Configuration
		oErr    string
	}{
		// "failed ping": {
		// 	iConfig: mysql.Configuration{},
		// },
		"mysql": {
			iConfig: mysql.Configuration{
				Hostname:       hostname,
				Port:           port,
				Username:       username,
				Password:       password,
				Database:       TestDatabaseName,
				ConnectTimeout: 10 * time.Second,
				ParseTime:      false,
			},
		},
	}
	for cDesc, c := range cases {
		database := mysql.NewMetaMySQL()
		//connect to databsae
		err := database.Initialize(c.iConfig)
		if c.oErr != "" {
			assert.NotNil(t, err, testCaseMap, cDesc)
			assert.Equal(t, c.oErr, err.Error(), testCaseMap, cDesc)
		} else {
			assert.Nil(t, err, testCaseMap, cDesc)
		}
		//shutdown from database
		err = database.Shutdown()
		//check error
		if c.oErr != "" {
			assert.NotNil(t, err, testCaseMap, cDesc)
			assert.Equal(t, c.oErr, err.Error(), testCaseMap, cDesc)
		} else {
			assert.Nil(t, err, testCaseMap, cDesc)
		}
	}
}

// func TestIntVerifyTable(t *testing.T) {
// 	//Test: the goal of this test is to create situations and determine whether or not the verify
// 	// function works as expected
// 	//Notes: We will need to have some way to create the tables
// 	//Verification: Verify that the table is found or not found

// 	configs := map[string]database.Configuration{
// 		"mysql; transactions": {
// 			Driver:          "mysql",
// 			Hostname:        mysqlHostname,
// 			UseTransactions: true},
// 		"mysql; no transactions": {
// 			Driver:   "mysql",
// 			Hostname: mysqlHostname}} //,
// 	// "sqlite; transactions": {
// 	// 	Driver:          "sqlite",
// 	// 	FilePath:        TestDatabaseFile,
// 	// 	UseTransactions: true},
// 	// "sqlite; no transactions": {
// 	// 	Driver:   "sqlite",
// 	// 	FilePath: TestDatabaseFile}}
// 	// "postgres": {
// 	// 	Driver:   "postgres",
// 	// 	Hostname: postgresHostname,
// 	// 	Port:     postgresPort,
// 	// 	Username: postgresRootUsername,
// 	// 	Password: rootPassword},

// 	cases := map[string]struct {
// 		iTable string //table to verify
// 		oFound bool   //whether or not the table should be found
// 		oErr   string //expected error
// 	}{
// 		"employee; found": {
// 			iTable: database.TableEmployee,
// 			oFound: true},
// 		"employee; not found": {
// 			iTable: database.TableEmployee,
// 			oFound: false},
// 		"client; found": {
// 			iTable: database.TableClient,
// 			oFound: true},
// 		"client; not found": {
// 			iTable: database.TableClient,
// 			oFound: false},
// 		"project; found": {
// 			iTable: database.TableProject,
// 			oFound: true},
// 		"project; not found": {
// 			iTable: database.TableProject,
// 			oFound: false},
// 		"timer; found": {
// 			iTable: database.TableTimer,
// 			oFound: true},
// 		"timer; not found": {
// 			iTable: database.TableTimer,
// 			oFound: false}}

// 	//create database struct
// 	database := database.NewDatabase()
// 	for confDesc, config := range configs {
// 		//connect to database
// 		if err := database.Connect(config); err != nil {
// 			t.Fatalf(testUnexpectedError+testConfigMap, t.Name(), err.Error(), confDesc)
// 		}
// 		for cDesc, c := range cases {
// 			//switch on test
// 			switch c.oFound {
// 			case false:
// 				switch c.iTable {
// 				case database.TableTimer:
// 					if err := database.DropTableTimer(); err != nil {
// 						t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 					}
// 				case database.TableProject:
// 					if err := database.DropTableProject(); err != nil {
// 						t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 					}
// 				case database.TableClient:
// 					if err := database.DropTableProject(); err != nil {
// 						t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 					}
// 					if err := database.DropTableClient(); err != nil {
// 						t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 					}
// 				case database.TableEmployee:
// 					// if err := database.DropTableTimer(); err != nil {
// 					// 	t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 					// }
// 					if err := database.DropTableEmployee(); err != nil {
// 						t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 					}
// 				default:
// 					t.Fatalf("%s, unsupported table configured: \"%s\", case: %s, config: %s", t.Name(), c.iTable, cDesc, confDesc)
// 				}
// 			case true:
// 				switch c.iTable {
// 				case database.TableTimer:
// 					if err := database.CreateTableEmployee(); err != nil {
// 						t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 					}
// 					if err := database.CreateTableTimer(); err != nil {
// 						t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 					}
// 				case database.TableProject:
// 					if err := database.CreateTableClient(); err != nil {
// 						t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 					}
// 					if err := database.CreateTableProject(); err != nil {
// 						t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 					}
// 				case database.TableClient:
// 					if err := database.CreateTableClient(); err != nil {
// 						t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 					}
// 				case database.TableEmployee:
// 					if err := database.CreateTableEmployee(); err != nil {
// 						t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 					}
// 				default:
// 					t.Fatalf("%s, unsupported table configured: \"%s\", case: %s, config: %s", t.Name(), c.iTable, cDesc, confDesc)
// 				}
// 			}
// 			//verify and check error
// 			if tables, err := database.VerifyTables(); err != nil {
// 				t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 			} else {
// 				if found, ok := tables[c.iTable]; ok {
// 					if found != c.oFound {
// 						t.Fatalf("TestIntVerifyTables, found doesn't match, got %t, expected %t, case: %s", found, c.oFound, cDesc)
// 					}
// 				} else {
// 					t.Fatalf("TestIntVerifyTables, table: \"%s\" unexpectedly not found, case: %s", c.iTable, cDesc)
// 				}
// 			}
// 		}
// 		//drop tables
// 		// database.DropTableTimer()
// 		// database.DropTableProject()
// 		// database.DropTableClient()
// 		database.DropTableEmployee()
// 		//disconnect from database
// 		if err := database.Disconnect(); err != nil {
// 			t.Fatalf(testUnexpectedError+"case: %s", t.Name(), err.Error(), confDesc)
// 		}
// 	}
// 	database.Close()
// 	database = nil
// }

// func TestIntEmployeeRead(t *testing.T) {
// 	//Test: the goal of this is to test the application's ability to successfully (and unsuccessfully) read
// 	// an employee
// 	//Notes: Employees don't have any foreign key constraints, so there shouldn't be much of an issue setting
// 	// up the database to work as expected
// 	//Verification: Verify that the employee read has the expected ID and employee values

// 	configs := map[string]database.Configuration{
// 		"mysql; transactions": {
// 			Driver:          "mysql",
// 			Hostname:        mysqlHostname,
// 			UseTransactions: true},
// 		"mysql; no transactions": {
// 			Driver:          "mysql",
// 			Hostname:        mysqlHostname,
// 			UseTransactions: false}} //,

// 	employees := map[int]bludgeon.Employee{
// 		1: bludgeon.Employee{
// 			ID:        1,
// 			FirstName: "Antonio",
// 			LastName:  "Alexander"},
// 		2: bludgeon.Employee{
// 			ID:        2,
// 			FirstName: "Reasonable",
// 			LastName:  "Doubt"},
// 		3: bludgeon.Employee{
// 			ID:        3,
// 			FirstName: "Antonio",
// 			LastName:  "Alexander"}}

// 	cases := map[string]struct {
// 		iID       int64
// 		oErr      string
// 		oEmployee bludgeon.Employee
// 	}{
// 		"normal; id 1": {
// 			iID:       1,
// 			oEmployee: employees[1]},
// 		"normal; id 2": {
// 			iID:       2,
// 			oEmployee: employees[2]},
// 		"normal; id 3": {
// 			iID:       3,
// 			oEmployee: employees[3]},
// 		"non-existant; id 4": {
// 			iID:  4,
// 			oErr: sql.ErrNoRows.Error()}}

// 	for confDesc, config := range configs {
// 		//create database
// 		database := database.NewDatabase()
// 		//connect to database
// 		if err := database.Connect(config); err != nil {
// 			t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), confDesc)
// 		}
// 		//drop employee table
// 		if err := database.DropTableEmployee(); err != nil {
// 			t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), confDesc)
// 		}
// 		//create employee table
// 		if err := database.CreateTableEmployee(); err != nil {
// 			t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), confDesc)
// 		}
// 		//create employees
// 		for _, employee := range employees {
// 			if _, err := database.EmployeeCreate(employee); err != nil {
// 				t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), confDesc)
// 			}
// 		}
// 		for cDesc, c := range cases {
// 			//read employee and check for error
// 			if employee, err := database.EmployeeRead(c.iID); err != nil {
// 				if err.Error() != c.oErr {
// 					t.Fatalf(testMismatchedError+testCaseMap, t.Name(), err.Error(), c.oErr, cDesc)
// 				}
// 			} else {
// 				//check id
// 				if employee.ID != c.oEmployee.ID {
// 					t.Fatalf(testMismatchedInteger+testCaseMap, t.Name(), "EmployeeID", employee.ID, c.oEmployee.ID, cDesc)
// 				}
// 				//check first name
// 				if employee.FirstName != c.oEmployee.FirstName {
// 					t.Fatalf(testMismatchedString+testCaseMap, t.Name(), "FirstName", employee.FirstName, c.oEmployee.FirstName, cDesc)
// 				}
// 				//check last name
// 				if employee.LastName != c.oEmployee.LastName {
// 					t.Fatalf(testMismatchedString+testCaseMap, t.Name(), "LastName", employee.LastName, c.oEmployee.LastName, cDesc)
// 				}
// 			}
// 		}
// 		//drop employee table
// 		if err := database.DropTableEmployee(); err != nil {
// 			t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), confDesc)
// 		}
// 		//disconnect from database
// 		if err := database.Disconnect(); err != nil {
// 			t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), confDesc)
// 		}
// 		//clean up database
// 		database = nil
// 	}
// }

// func TestIntEmployeesRead(t *testing.T) {
// 	//Test: the goal of this test is to confirm the APIs ability to read multiple employees
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql; transactions": {
// 			Driver:          "mysql",
// 			Hostname:        mysqlHostname,
// 			UseTransactions: true},
// 		"mysql; no transactions": {
// 			Driver:          "mysql",
// 			Hostname:        mysqlHostname,
// 			UseTransactions: false}} //,

// 	employees := map[int]bludgeon.Employee{
// 		1: bludgeon.Employee{
// 			ID:        1,
// 			FirstName: "Antonio",
// 			LastName:  "Alexander"},
// 		2: bludgeon.Employee{
// 			ID:        2,
// 			FirstName: "Reasonable",
// 			LastName:  "Doubt"},
// 		3: bludgeon.Employee{
// 			ID:        3,
// 			FirstName: "Antonio",
// 			LastName:  "Alexander"}}

// 	cases := map[string]struct {
// 		iEmployees []bludgeon.Employee
// 		oErr       string
// 	}{
// 		"no employees": {
// 			oErr: sql.ErrNoRows.Error()},
// 		"single employee": {
// 			iEmployees: []bludgeon.Employee{
// 				employees[1]}},
// 		"multiple employes": {
// 			iEmployees: []bludgeon.Employee{
// 				employees[1],
// 				employees[2],
// 				employees[3]}}}

// 	for confDesc, config := range configs {
// 		//create database
// 		database := database.NewDatabase()
// 		//connect to database
// 		if err := database.Connect(config); err != nil {
// 			t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), confDesc)
// 		}
// 		for cDesc, c := range cases {
// 			//drop employee table
// 			if err := database.DropTableEmployee(); err != nil {
// 				t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), confDesc)
// 			}
// 			//create employee table
// 			if err := database.CreateTableEmployee(); err != nil {
// 				t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), confDesc)
// 			}
// 			//create employees
// 			for _, employee := range c.iEmployees {
// 				if _, err := database.EmployeeCreate(employee); err != nil {
// 					t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), confDesc)
// 				}
// 			}
// 			//read employees and check for error
// 			if employees, err := database.EmployeesRead(); err != nil {
// 				if err.Error() != c.oErr {
// 					t.Fatalf(testMismatchedError+testCaseMap, t.Name(), err.Error(), c.oErr, cDesc)
// 				}
// 			} else {
// 				//check for length
// 				if len(employees) != len(c.iEmployees) {
// 					t.Fatalf(testMismatchedInteger+testCaseMap, t.Name(), "EmployeeLength", len(employees), len(c.iEmployees), cDesc)
// 				}
// 				//check employees
// 				for i, employee := range c.iEmployees {
// 					//check id
// 					if employee.ID != c.iEmployees[i].ID {
// 						t.Fatalf(testMismatchedInteger+testCaseMap, t.Name(), "EmployeeID", employee.ID, c.iEmployees[i].ID, cDesc)
// 					}
// 					//check first name
// 					if employee.FirstName != c.iEmployees[i].FirstName {
// 						t.Fatalf(testMismatchedString+testCaseMap, t.Name(), "FirstName", employee.FirstName, c.iEmployees[i].FirstName, cDesc)
// 					}
// 					//check last name
// 					if employee.LastName != c.iEmployees[i].LastName {
// 						t.Fatalf(testMismatchedString+testCaseMap, t.Name(), "LastName", employee.LastName, c.iEmployees[i].LastName, cDesc)
// 					}
// 				}
// 			}
// 			//drop employee table
// 			if err := database.DropTableEmployee(); err != nil {
// 				t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), confDesc)
// 			}
// 		}
// 		//disconnect from database
// 		if err := database.Disconnect(); err != nil {
// 			t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), confDesc)
// 		}
// 		//clean up database
// 		database = nil
// 	}
// }

// func TestIntEmployeesUpdate(t *testing.T) {
// 	//Test: the goal of the test is to test the api's ability to update an employee
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql; transactions": {
// 			Driver:          "mysql",
// 			Hostname:        mysqlHostname,
// 			UseTransactions: true},
// 		"mysql; no transactions": {
// 			Driver:          "mysql",
// 			Hostname:        mysqlHostname,
// 			UseTransactions: false}}

// 	employees := map[int]bludgeon.Employee{
// 		1: bludgeon.Employee{
// 			FirstName: "Antonio",
// 			LastName:  "Alexander"}}

// 	cases := map[string]struct {
// 		iEmployeeCreate bludgeon.Employee //employee to initially create
// 		iEmployeeUpdate bludgeon.Employee //employee to update
// 		oErr            string            //error returned after updating
// 	}{
// 		"normal employee": {
// 			iEmployeeCreate: employees[1],
// 			iEmployeeUpdate: bludgeon.Employee{
// 				FirstName: "Tony"}},
// 		"non-existant employee": {
// 			iEmployeeUpdate: bludgeon.Employee{
// 				ID:        -1,
// 				FirstName: "Unreasonable"},
// 			oErr: database.ErrUpdateFailed}}
// 	//attempt to update an employee incorrectly?

// 	//create database
// 	database := database.NewDatabase()
// 	for confDesc, config := range configs {
// 		//connect to database
// 		if err := database.Connect(config); err != nil {
// 			t.Fatalf(testUnexpectedError+testConfigMap, t.Name(), err.Error(), confDesc)
// 		}
// 		for cDesc, c := range cases {
// 			//drop employee table
// 			if err := database.DropTableEmployee(); err != nil {
// 				t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 			}
// 			//create employee table
// 			if err := database.CreateTableEmployee(); err != nil {
// 				t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 			}
// 			//create employee (don't check for errors)
// 			if id, err := database.EmployeeCreate(c.iEmployeeCreate); err != nil {
// 				t.Logf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 			} else {
// 				//read employee and verify write
// 				if employee, err := database.EmployeeRead(id); err != nil {
// 					t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 				} else {
// 					if employee.FirstName != c.iEmployeeCreate.FirstName {
// 						t.Fatalf(testMismatchedString+testCaseMap+testConfigMap, t.Name(), "FirstName", employee.FirstName, c.iEmployeeCreate.FirstName, cDesc, confDesc)
// 					}
// 					if employee.LastName != c.iEmployeeCreate.LastName {
// 						t.Fatalf(testMismatchedString+testCaseMap+testConfigMap, t.Name(), "FirstName", employee.LastName, c.iEmployeeCreate.LastName, cDesc, confDesc)
// 					}
// 				}
// 				//update employee and check for errors
// 				if err := database.EmployeeUpdate(id, c.iEmployeeUpdate); err != nil {
// 					if err.Error() != c.oErr {
// 						t.Fatalf(testMismatchedError+testCaseMap+testConfigMap, t.Name(), err.Error(), c.oErr, cDesc, confDesc)
// 					}
// 				} else {
// 					//read employee and verify changes
// 					if employee, err := database.EmployeeRead(id); err != nil {
// 						t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 					} else {
// 						if employee.FirstName != c.iEmployeeUpdate.FirstName {
// 							t.Fatalf(testMismatchedString+testCaseMap+testConfigMap, t.Name(), "FirstName", employee.FirstName, c.iEmployeeUpdate.FirstName, cDesc, confDesc)
// 						}
// 						if employee.LastName != c.iEmployeeUpdate.LastName {
// 							t.Fatalf(testMismatchedString+testCaseMap+testConfigMap, t.Name(), "FirstName", employee.LastName, c.iEmployeeUpdate.LastName, cDesc, confDesc)
// 						}
// 					}
// 				}
// 			}
// 			//drop employee table
// 			if err := database.DropTableEmployee(); err != nil {
// 				t.Fatalf(testUnexpectedError+testCaseMap+testConfigMap, t.Name(), err.Error(), cDesc, confDesc)
// 			}
// 		}
// 		//disconnect from database
// 		if err := database.Disconnect(); err != nil {
// 			t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), confDesc)
// 		}
// 	}
// 	//clean up database
// 	database = nil
// }

// func TestIntEmployeeDelete(t *testing.T) {
// 	//Test: The goal of this test is just to verify the ability to delete an employee
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql; transactions": {
// 			Driver:          "mysql",
// 			Hostname:        mysqlHostname,
// 			UseTransactions: true},
// 		"mysql; no transactions": {
// 			Driver:          "mysql",
// 			Hostname:        mysqlHostname,
// 			UseTransactions: false}} //,

// 	cases := map[string]struct {
// 		iEmployeeCreate   bludgeon.Employee //employee to create
// 		iEmployeeDeleteID int64             //employee to delete
// 		oErr              string            //expected error for delete
// 	}{}
// 	//attempt to delete existing employee
// 	//attempt to delete an employee that doesn't exist

// 	//create database
// 	database := database.NewDatabase()
// 	for confDesc, config := range configs {
// 		//connect to database
// 		if err := database.Connect(config); err != nil {
// 			t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), confDesc)
// 		}
// 		for cDesc, c := range cases {
// 			//drop employee table
// 			if err := database.DropTableEmployee(); err != nil {

// 			}
// 			//create employee table
// 			if err := database.CreateTableEmployee(); err != nil {

// 			}
// 			//create employee (don't check for errors)
// 			if id, err := database.EmployeeCreate(c.iEmployeeCreate); err != nil {
// 				t.Logf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), cDesc)
// 			} else {
// 				//read employee and verify write
// 				if employee, err := database.EmployeeRead(id); err != nil {
// 					t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), cDesc)
// 				} else {
// 					if employee.FirstName != c.iEmployeeCreate.FirstName {
// 						t.Fatalf(testMismatchedString+testCaseMap, t.Name(), "FirstName", employee.FirstName, c.iEmployeeCreate.FirstName, cDesc)
// 					}
// 					if employee.LastName != c.iEmployeeCreate.LastName {
// 						t.Fatalf(testMismatchedString+testCaseMap, t.Name(), "FirstName", employee.LastName, c.iEmployeeCreate.LastName, cDesc)
// 					}
// 				}
// 			}
// 			//delete employee
// 			if err := database.EmployeeDelete(c.iEmployeeDeleteID); err != nil {
// 				if err.Error() != c.oErr {
// 					t.Fatalf(testMismatchedError+testCaseMap, t.Name(), err.Error(), c.oErr, cDesc)
// 				}
// 			} else {
// 				if _, err := database.EmployeeRead(c.iEmployeeDeleteID); err != nil {
// 					t.Fatalf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), cDesc)
// 				}
// 			}
// 			//read employee
// 			//drop employee table
// 			if err := database.DropTableEmployee(); err != nil {
// 				t.Logf(testUnexpectedError+testCaseMap, t.Name(), err.Error(), cDesc)
// 			}
// 		}
// 		//disconnect from database
// 		database.Disconnect()
// 	}
// 	//destroy database
// 	database = nil
// }

// func TestIntTimerRead(t *testing.T) {
// 	//Test:
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql": {
// 			Driver:   "mysql",
// 			Hostname: mysqlHostname}}

// 	cases := map[string]struct{}{}
// 	for range configs {
// 		for range cases {

// 		}
// 	}
// }
// func TestIntTimersRead(t *testing.T) {
// 	//Test:
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql": {
// 			Driver:   "mysql",
// 			Hostname: mysqlHostname}}

// 	cases := map[string]struct{}{}
// 	for range configs {
// 		for range cases {

// 		}
// 	}
// }

// func TestIntTimerUpdate(t *testing.T) {
// 	//Test:
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql": {
// 			Driver:   "mysql",
// 			Hostname: mysqlHostname}}

// 	cases := map[string]struct{}{}
// 	for range configs {
// 		for range cases {

// 		}
// 	}
// }

// func TestIntTimerDelete(t *testing.T) {
// 	//Test:
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql": {
// 			Driver:   "mysql",
// 			Hostname: mysqlHostname}}

// 	cases := map[string]struct{}{}
// 	for range configs {
// 		for range cases {

// 		}
// 	}
// }

// func TestIntClientRead(t *testing.T) {
// 	//Test:
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql": {
// 			Driver:   "mysql",
// 			Hostname: mysqlHostname}}

// 	cases := map[string]struct{}{}
// 	for range configs {
// 		for range cases {

// 		}
// 	}
// }

// func TestIntClientsRead(t *testing.T) {
// 	//Test:
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql": {
// 			Driver:   "mysql",
// 			Hostname: mysqlHostname}}

// 	cases := map[string]struct{}{}
// 	for range configs {
// 		for range cases {

// 		}
// 	}
// }

// func TestIntClientUpdate(t *testing.T) {
// 	//Test:
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql": {
// 			Driver:   "mysql",
// 			Hostname: mysqlHostname}}

// 	cases := map[string]struct{}{}
// 	for range configs {
// 		for range cases {

// 		}
// 	}
// }
// func TestIntClientDelete(t *testing.T) {
// 	//Test:
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql": {
// 			Driver:   "mysql",
// 			Hostname: mysqlHostname}}

// 	cases := map[string]struct{}{}
// 	for range configs {
// 		for range cases {

// 		}
// 	}
// }

// func TestIntProjectRead(t *testing.T) {
// 	//Test:
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql": {
// 			Driver:   "mysql",
// 			Hostname: mysqlHostname}}

// 	cases := map[string]struct{}{}
// 	for range configs {
// 		for range cases {

// 		}
// 	}
// }

// func TestIntProjectsRead(t *testing.T) {
// 	//Test:
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql": {
// 			Driver:   "mysql",
// 			Hostname: mysqlHostname}}

// 	cases := map[string]struct{}{}
// 	for range configs {
// 		for range cases {

// 		}
// 	}
// }

// func TestIntProjectUpdate(t *testing.T) {
// 	//Test:
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql": {
// 			Driver:   "mysql",
// 			Hostname: mysqlHostname}}

// 	cases := map[string]struct{}{}
// 	for range configs {
// 		for range cases {

// 		}
// 	}
// }

// func TestIntProjectDelete(t *testing.T) {
// 	//Test:
// 	//Notes:
// 	//Verification:

// 	configs := map[string]database.Configuration{
// 		"mysql": {
// 			Driver:   "mysql",
// 			Hostname: mysqlHostname}}

// 	cases := map[string]struct{}{}
// 	for range configs {
// 		for range cases {

// 		}
// 	}
// }
