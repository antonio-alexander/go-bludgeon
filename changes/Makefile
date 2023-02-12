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

# Reference: https://medium.com/@pedram.esmaeeli/generate-swagger-specification-from-go-source-code-648615f7b9d9
check-swagger: ## - validate/install swagger (v0.29.0)
	which swagger || (go install github.com/go-swagger/go-swagger/cmd/swagger@v0.29.0)

swagger: check-swagger ## - generate the swagger.json
	@swagger generate spec --work-dir=./internal/swagger --output ./tmp/swagger.json --scan-models

validate-swagger: swagger ## - validate the swagger.json
	@swagger validate ./tmp/swagger.json

serve-swagger: swagger ## - serve (web) the swagger.json
	@swagger serve -F=swagger ./tmp/swagger.json -p ${swagger_port} --no-open

check-godoc: ## - validate/install godoc
	which godoc || (go install golang.org/x/tools/cmd/godoc@v0.1.10)

serve-godoc: check-godoc ## - serve (web) the godocs
	@cd .. && godoc -http :8080

test: run ## - test the source
	@go test -v -parallel=1 -cover -count=1 ./... -coverprofile ./tmp/go-bludgeon-changes.out | tee ./tmp/go-bludgeon-changes.log

build: ## - build the source (latest)
	@docker compose --profile application build --build-arg GIT_COMMIT=`git rev-parse HEAD` --build-arg GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
	@docker image prune -f

dep: ## - run only the dependencies (docker) detached
	@docker compose up -d --wait

run: ## - run the service and its dependencies (docker)
	@docker container rm -f changes
	@docker image prune -f
	@docker compose --profile application up -d --wait

stop: ## - stop and clean docker resources
	@docker compose --profile application down
	@docker compose --profile application rm -f
