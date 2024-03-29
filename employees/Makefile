## ----------------------------------------------------------------------
## This makefile can be used to execute common functions to interact with
## the source code, these functions ease local development and can also be
## used in CI/CD pipelines.
## ----------------------------------------------------------------------

swagger_port = 9000 # defaults

# REFERENCE: https://stackoverflow.com/questions/16931770/makefile4-missing-separator-stop
help: ## - Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

check-proto: ## - check protoc/proto-gen-go/protolint
	@which protolint || (go install github.com/yoheimuta/protolint/cmd/protolint@v0.38.3)
	@which protoc || echo protoc v3.20.1 not installed, install https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.1
	@which protoc-gen-go || (go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0)
	@which protoc-gen-go-grpc || (go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2)

lint-proto: check-proto ## - lint proto
	@protolint -fix ./data/pb/employees.proto

clean-proto: ## - clean protos
	@rm ./data/pb/employees.pb.go ./data/pb/employees_grpc.pb.go

build-proto: lint-proto ## - build proto
	@protoc -I="./data/pb" --go_opt=paths=source_relative --go_out="./data/pb" --go-grpc_opt=paths=source_relative --go-grpc_out="./data/pb" ./data/pb/employees.proto

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

test: run ## - test the source with verbose output
	@go test -v -cover -parallel=1 --count=1 ./... -coverprofile ./tmp/go-bludgeon-employees.out | tee ./tmp/go-bludgeon-employees.log

build: ## - build the source (latest)
	@docker compose --profile application build --build-arg GIT_COMMIT=`git rev-parse HEAD` --build-arg GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
	@docker image prune -f

dep: ## - run only the dependencies (docker) detached
	@docker compose up -d

run: ## - run the service and its dependencies (docker) detached
	@docker compose --profile application up -d

stop:
	@docker compose --profile application down

clean: stop ## - stop and clean docker resources
	@docker compose --profile application rm -f