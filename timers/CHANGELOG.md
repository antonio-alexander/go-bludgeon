# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.2] - 2023-02-22

- fixed security vulnerabilities by updating volumes
- upgraded to golang.org/x/text v0.3.8
- upgraded to golang.org/x/net v0.7.0

## [1.3.1] - 2023-02-14

- updated the changes client
- updated logic for changes to use go routines

## [1.3.0] - 2023-02-12

- updated to latest internal
- integrated changes client

## [1.1.3] - 2022-07-14

- fixed bug with swagger document search parameters
- added grpc functionality

## [1.1.2] - 2022-06-25

- Updated documentation
- Added missing endpoints on rest service (archive/comment)

## [1.1.1] - 2022-06-05

- Removed Audit type, copied contents to Timers and TimeSlices types
- Resolved broken code in tests/meta
- Updated service to allow operations with time slices
- Updated client to include endpoints for time slices
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
