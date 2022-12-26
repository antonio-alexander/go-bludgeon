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
	@golangci-lint run --verbose

dep: ## - run only the dependencies (docker) detached
	@docker compose up -d

test: dep ## - test the source with verbose output
	@go test -v -cover -parallel=1 --count=1 ./... -coverprofile ./tmp/go-bludgeon-internal.out | tee ./tmp/go-bludgeon-internal.log

