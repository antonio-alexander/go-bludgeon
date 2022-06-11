## ----------------------------------------------------------------------
## This makefile can be used to execute common functions to interact with
## the source code, these functions ease local development and can also be
## used in CI/CD pipelines.
## ----------------------------------------------------------------------

# REFERENCE: https://stackoverflow.com/questions/16931770/makefile4-missing-separator-stop
help: ## - Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

check-lint: ## - validate/install golangci-lint installation
	which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2)

lint: check-lint ## - lint the source
	golangci-lint run

lint-verbose: check-lint ## - lint the source with verbose output
	golangci-lint run --verbose

# Reference: https://medium.com/@pedram.esmaeeli/generate-swagger-specification-from-go-source-code-648615f7b9d9
check-swagger: ## - validate/install swagger (v0.29.0)
	which swagger || (go install github.com/go-swagger/go-swagger/cmd/swagger@v0.29.0)

swagger: check-swagger ## - generate the swagger.json
	swagger generate spec --work-dir=./service/rest/swagger --output ./cmd/service/swagger.json --scan-models

validate-swagger: swagger ## - validate the swagger.json
	swagger validate ./cmd/service/swagger.json

serve-swagger: swagger ## - serve (web) the swagger.json
	swagger serve -F=swagger ./cmd/service/swagger.json -p 8080 --no-open

check-godoc: ## - validate/install godoc
	which godoc || (go install golang.org/x/tools/cmd/godoc@v0.1.10)

serve-godoc: check-godoc ## - serve (web) the godocs
	cd .. && godoc -http :8080

test: ## - test the source
	 go test -cover --count=1 ./...

test-verbose: ## - test the source with verbose output
	 go test -cover -v --count=1 ./...

build: ## - build the source (latest)
	./build.sh latest

run: ## - run the service and its dependencies (docker)
	docker compose up