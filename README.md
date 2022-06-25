# go-bludgeon (github.com/antonio-alexander/go-bludgeon)

go-bludgeon is a collection of microservices used to track time for one or more individuals. It's an "example" repo for microservices to serve as a test bed for common things like github workflows, git flow, documentation etc. It's high-level purpose is to provide a kind of battle-tested/lessons-learned architecture in pure Go, that doesn't necessarily make opinionated decisions, but makes decisions with a specific idea in mind.

## Getting Started

To get started, use the [docker-compose.yml](./docker-compose.yml) to pull the most recent images and to get everything up and running. The Swagger and Godocs images can be used to read the most recent documentation and swagger can execute each endpoint with examples.

> Keep in mind that this docker-compose CAN'T be used to build anything.

## MySQL

The MySQL image contains all of sql required for all of the mysql/database, its used when mysql is enabled as a data store, for more information look at the mysql [README.md](./mysql/README.md)

## Employees

The employees service is used to manage the employees object, for more information look at the employees [README.md](./employees/README.md).

## Timers

The timers service is used to manage objects related to timers, for more information look at the employees [README.md](./timers/README.md).
