# Changes (github.com/antonio-alexander/go-bludgeon/changes)

Changes is a service that can be used to broadcast when a change has occured in a given service with a specific kind of data. This can be used to communicate changes between services in a way that maintains the disconnection between services. This is for situations where microservices with different concerns have no authority with certain objects but depends on them.

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
