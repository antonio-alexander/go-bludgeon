# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.4.3] - 2023-02-22

- upgraded to golang.org/x/net v0.7.0

## [1.4.2] - 2023-02-22

- fixed security vulnerabilities by updating volumes
- upgraded to golang.org/x/text v0.3.8
- upgraded from Go 1.16 to 1.19

## [1.4.1] - 2023-02-12

- refactored some of the logic in the websockets client

## [1.4.0] - 2023-01-21

- updated errors types to be a little more friendly for marshal/unmarshal
- refactored the rest client to expose the status code (and not just bytes/error)
- updated the rest/grpc service to use their parameters to register the endpoints/services

## [1.3.2] - 2022-12-26

- Had an issue with caching (go-proxy); had to update the version to create a new [valid] tag so go mod downloads would work; no functional changes to the code were made

## [1.3.0] - 2022-12-06

- Added Initializer/Configurer/Shutdown/Closer interfaces
- Updated internal packages to use common interfaces

## [1.2.0] - 2022-07-27

- Added context for client
- Added GRPC client and server

## [1.0.0] - 2021-03-27

- Initial release
