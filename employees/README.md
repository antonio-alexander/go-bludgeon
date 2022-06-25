# Employees (github.com/antonio-alexander/go-bludgeon/employees)

Employees is a service that can be used to interact with employees (only employee). An employee is simply a resource that can be used to describe "someone" within the context of go-bludgeon.

Employees is pretty vanilla/boring. It allows CRUD (Create, Read, Update, Delete) operations on a given employee object.

The service architecture is layered and fairly common, the meta is wrapped by the logic which is then exposed by the service. Multiple options are provided for meta (memory, file and database). Configuration can be used to configure which to use and the manner of its usage. The available implementations for meta are:

- Memory: the scope of this implementation is a single application instance, it provides no communication between multiple instances
- File: the scope of this implementation is a file system, it allows configuration of file locking such that multiple instances using the same file system can work concurrently
- Database: the scope of this impelemtnation is a database, it can be used concurrently

## Getting Started

The service (and its dependencies) can be easily started by using the [docker-compose.yml](./docker-compose.yml) file, to build and run everything. The docker-compose can be used to build the employees service (but you can also pull the image).

Use this command to bring up the service and its dependencies:

```sh
docker compose up -d
```

use one of these comands to build the service (see the [build.sh](./build.sh)):

```sh
docker compose build
```

```sh
./build.sh latest
```

> Keep in mind that the build.sh is a better solution to build locally since it'll embed git information into the image itself

## Dependencies

This service has a single dependency (if the MySQL meta is enabled). It is not a hard dependency, but it will limit what you can do with the service. Although it is possible to create the mysql tables on your own, the mysql image makes this much easier. See the [docker-compose.yml](./docker-compose.yml) for more information.

## Configuration

The employees service has some ability to configure itself along with the current feature set, expect that as this feature set changes, the configuration will also change. The service can be configured using environmental variables or with a configuration file (to be loaded via volume).

### Meta

The following environmental variables can be used to configure meta:

- BLUDGEON_META_TYPE: this can be used to configure the implementation of meta
  - valid values: mysql, memory, file
- DATABASE_HOST: this is the address of the database
  - example: mysql, host.docker.internal
- DATABASE_PORT: this is the port of the database
  - examples: 3306
- DATABASE_NAME: this is the name of the database
  - example values: bludgeon
- DATABASE_USER: this is the user to use when connecting to the database
  - example values: bludgeon
- DATABASE_PASSWORD: this is the password to use when connecting to the database
  - example: bludgeon
- BLUDGEON_META_FILE: this is the file to use when using the file meta
  - example: "data/bludgeon.json"

### Server

The following environmental variables can be used to configure the server:

- BLUDGEON_SERVICE_TYPE: this configures the type of service running, only rest is supported currently
  - valid values: rest
- BLUDGEON_REST_ADDRESS: this is the address the server listens on
  - valid values: ""
- BLUDGEON_REST_PORT: this is the port the server listens on
  - example values: 8080
- BLUDGEON_ALLOWED_ORIGINS: this is rest specific and configures which origins are allowed (CORS)
  - example values: "http://host.docker.internal"
- BLUDGEON_CORS_DEBUG: this configures whether or not CORS debug is configured
  - valid values: true, false

## Contracts

These are contracts used to communicate with the employees service.

### Employee

Employee is a representation of an "employee" when read. Keep in mind that an employee's natural key or rather what makes it "unique" is the email address.

```json
{
    "id": "86fa2f09-d260-11ec-bd5d-0242c0a8e002",
    "first_name": "John",
    "last_name": "Smith",
    "email_address": "John.Smith@foobar.duck",
    "last_updated": 1652417242000,
    "last_updated_by": "bludgeon_employee_memory",
    "version": 1
}
```

### Employee Partial

Employee Partial is a representation of an "employee" when it must be mutated. Note the absence of audit fields such as last_updated, last_updated_by and version as well as the absence of id. When mutating an employee, this will be the contract used.

```json
{
    "first_name": "Jane",
    "last_name": "Doe",
    "email_address": "Jane.Doe@foobar.duck",
}
```

### Employee Search

Employee Search is used to read employees, it can be used to filter (or search) for empployees depending on the values in the properties.

```json
{
    "ids": ["5afbea80-f36e-4e20-8763-64b5badbaf7d","a059dd1e-8406-4f8e-ace0-6cd2a5d54166"],
    "first_name": "Jane",
    "first_names": ["John", "Jane"],
    "last_name": "Doe",
    "last_names":  ["John", "Jane"],
    "email_address": "name@company.com",
    "email_addresses": ["john.doe@company.com", "jane.doe@company.com"],
}
```

## Business Logic

This business logic is exposed by the service

- Create employee: this will allow you to create an employee using the [employee partial contract](#employee-partial); this is the only mutation endpoint that recognizes email address; it cannot be changed post create; it will return a scalar employee using the [employee contract](#employee)
- Read employee: this will allow you to read an employee by providing its id; it will return a scalar employee using the [employee contract](#employee)
- Search for employees: this will allow you to read zero or more employees by providing information wiithin the [employee search contract](#employee-search); it will return an array of employees using the [employee contract](#employee)
- Update an employee: this will allow you to mutate an existing employee by providing the id and the [employee partial contract](#employee-partial); it will return a scalar employee using the [employee contract](#employee)
- Delete an employee: this will allow you to delete an employee that exists using it's id
