# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.1] - 2023-01-07

- updated rest client to use latest internal (for rest)
- normalized types for rest/grpc client so they're inter-changeable
- fixed service configuration to expose grpc (config issue)

## [1.3.0] - 2022-12-26

- integrated changes client

## [1.2.0] - 2022-07-30

- added grpc service and client
- updated internal package

## [1.1.3] - 2022-07-14

- fixed bug with swagger document search parameters

## [1.1.2] - 2022-06-25

- Updated documentation

## [1.1.1] - 2022-06-05

- Removed Audit type, copied contents to Timers and TimeSlices types
- Resolved broken code in tests/meta
- Updated github actions to validate swagger

## [1.1.0] - 2022-05-28

- Updated to microservice architecture
- Added swagger for rest
- Added golangci-lint
- Added Makefile
- Added client
- Added tests

## [1.0.0] - 2021-03-27

- Initial release
