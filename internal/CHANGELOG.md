# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.1] - 2022-12-26

- Had an issue with caching (go-proxy); had to update the version to create a new [valid] tag so go mod downloads would work; no functional changes to the code were made

## [1.3.0] - 2022-12-06

- Added Initializer/Configurer/Shutdown/Closer interfaces
- Updated internal packages to use common interfaces

## [1.2.0] - 2022-07-27

- Added context for client
- Added GRPC client and server

## [1.0.0] - 2021-03-27

- Initial release
