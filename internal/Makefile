## ----------------------------------------------------------------------
## This makefile can be used to execute common functions to interact with
## the source code, these functions ease local development and can also be
## used in CI/CD pipelines.
## ----------------------------------------------------------------------

swagger_port = 8082 # defaults

# REFERENCE: https://stackoverflow.com/questions/16931770/makefile4-missing-separator-stop
help: ## - Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

check-lint: ## - validate/install golangci-lint installation
	which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2)

lint: check-lint ## - lint the source
	@golangci-lint run

lint-verbose: check-lint ## - lint the source with verbose output
	@golangci-lint run --verbose

test: ## - test the source
	@go test -cover -p 1 --count=1 ./...

test-verbose: ## - test the source with verbose output
	@go test -cover -v -p 1 --count=1 ./...
