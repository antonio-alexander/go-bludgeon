# Timers (github.com/antonio-alexander/go-bludgeon/timers)

Timers is a service that can be used to interact with timers objects such as a timer and a time slice. The service allows you to mutate, start, stop and read timer objects. A timer can be used to track time spent on a given task. A timer can also be related to one (or none) employees (managed by the [employees service](../employees/)).

This implementation of timers is a bit unique in that behind the scenes time isn't simply the difference between start and stop time, but the sum of the elapsed time of a time slice. A time slice is the primitive unit which is like a timer, but the difference is that the elapsed time **actually** is the difference between finish and start time.

The service architecture is layered and fairly common, the meta is wrapped by the logic which is then exposed by the service. Multiple options are provided for meta (memory, file and database). Configuration can be used to configure which to use and the manner of its usage. The available implementations for meta are:

- Memory: the scope of this implementation is a single application instance, it provides no communication between multiple instances
- File: the scope of this implementation is a file system, it allows configuration of file locking such that multiple instances using the same file system can work concurrently
- Database: the scope of this impelemtnation is a database, it can be used concurrently

## Getting Started

The service (and its dependencies) can be easily started by using the [docker-compose.yml](./docker-compose.yml) file, to build and run everything. The docker-compose can be used to build the timers service (but you can also pull the image).

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

This service has two main dependencies: employees and mysql. They are not hard dependencies, but the absence of them will limit what you can do with the service. Employees can be used to get valid IDs for configuring the employee_id of a given timer and mysql is necessary if you don't want to have to deploy the mysql tables on your own. See the [docker-compose.yml](./docker-compose.yml) for more information.

## Configuration

The timers service has some ability to configure itself along with the current feature set, expect that as this feature set changes, the configuration will also change. The service can be configured using environmental variables or with a configuration file (to be loaded via volume).

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

### Timer

Timer is a representation of a "timer" when read. A timer doesn't really have a natural key, the only thing that makes it "unique" is the id itself.

```json
{
    "completed": false,
    "archived": false,
    "start": 1653719208,
    "finish": 1653719229,
    "elapsed_time": 21,
    "employee_id": "2e3a4156-b415-4120-982f-399182e99588",
    "active_time_slice_id":"a33f813e-e9bc-46ad-9956-0c4b6c1367ab",
    "id":"24dfe1eb-26a7-41db-a647-fe6cc5e77ab8",
    "comment":"This is a timer for lunch",
    "last_updated": 1652417242000,
    "last_updated_by": "bludgeon_employee_memory",
    "version": 1
}
```

### Timer Partial

Timer partial is a representation of a "timer" when it must be mutated. Note the absence of audit fields such as last_updated, last_updated_by and version as well as the absence of id. When mutating a timer, this is the contract used.

```json
{
    "completed": false,
    "archived": false,
    "employee_id": "2e3a4156-b415-4120-982f-399182e99588",
    "comment":"This is a timer for lunch",
    "finish": 1653719229,
}
```

### Timer Search

Timer search can be used to find one or more timers using a number of search parameters.

```json
{
    "employee_id": "2e3a4156-b415-4120-982f-399182e99588",
    "employee_ids": "2e3a4156-b415-4120-982f-399182e99588, 2e3a4156-b415-4120-982f-399182e99588",
    "completed": true,
    "archived": false,
    "ids":"24dfe1eb-26a7-41db-a647-fe6cc5e77ab8, 24dfe1eb-26a7-41db-a647-fe6cc5e77ab8"
}
```

### Time Slice

Time Slice is a representation of a time slice when read. It can be thought of as a "mini timer" that **can't** be paused. The two things that make a time slice "unique" is the id and the timer_id.

```json
{
    "completed": true,
    "start":1653720177,
    "finish": 1653720184,
    "elapsed_time": 7,
    "id": "ff7e87af-e6c5-44c3-851f-8801a33ad888",
    "timer_id": "7f583116-c7b8-457d-97e0-be0670e9e27e",
    "last_updated": 1652417242000,
    "last_updated_by": "bludgeon_employee_memory",
    "version": 1
}
```

### Time Slice Partial

Time Slice Partial is used when time slices need to be mutated. Note the absence of audit fields such as last_updated, last_updated_by and version as well as the absence of id.

```json
{
    "completed": true,
    "start":1653720177,
    "finish": 1653720184,
    "timer_id": "7f583116-c7b8-457d-97e0-be0670e9e27e"
}
```

### Time Slice Search

Time Slice Search can be used to read one or more time slices using a number of search parameters.

```json
{
    "employee_id": "2e3a4156-b415-4120-982f-399182e99588",
    "employee_ids": "2e3a4156-b415-4120-982f-399182e99588, 2e3a4156-b415-4120-982f-399182e99588",
    "completed": true,
    "archived": false,
    "ids":"24dfe1eb-26a7-41db-a647-fe6cc5e77ab8, 24dfe1eb-26a7-41db-a647-fe6cc5e77ab8"    
}
```

## Business logic

This business logic is exposed by the service

- Create a timer: this will allow you to create a timer using the [timer partial contract](#timer-partial)
- Read a timer: this will allow you to read a timer by providing the id; it will return a scalar timer using the [timer contract](#timer)
- Search for timers: this will allow you to search for zero or more timers using the [timer search contract](#timer-search); ; it will return an array of timers using the [timer contract](#timer)
- Delete a timer: this will allow you to delete a timer using the id
- Start a timer: this will allow you to start a timer using its id; it will return a scalar timer using the [timer contract](#timer)
- Stop a timer: this will allow you to stop a timer using its id; it will return a scalar timer using the [timer contract](#timer)
- Submit a timer: this will allow you to submit your timer using its id; it will return a scalar timer using the [timer contract](#timer)
- Update a timer comment: this will allow you to update the timers comment, the id and comment are required and the [timer partial contract](#timer-partial) will be used with only the comment value; it will return a scalar timer using the [timer contract](#timer)
- Archive Timer: this will allow you to archive the timer, the id and archive are required and the [timer partial contract](#timer-partial) will be used with only the archive value; it will return a scalar timer using the [timer contract](#timer)

Time slices aren't really meant to be used interactively, but are provided for administration. Under normal circumstances, these endpoints aren't expected to be used (at all); except by timer functions.

- Create a time slice: this will allow you to create a time slice using the [time slice partial contract](#time-slice-partial); this is the only endpoint that will allow you to set the timer id and returns a scalar time slice using the [time slice contract](#time-slice)
- Read a time slice: this will allow you to read a time slice using its id and will return a scalar time slice using the [time slice contract](#time-slice)
- Update a time slice: this can be used to update a time slice using the [time slice partial contract](#time-slice-partial) and will return a scalar time slice using the [time slice contract](#time-slice)
- Delete a time slice: this can be used to delete a time slice using its id
- Search for time slices: this can be used to search for zero or more time slices using the [time slice search contract](#time-slice-search) and will return an array of time slices using the [time slice contract](#time-slice)
