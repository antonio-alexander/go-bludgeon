## ----------------------------------------------------------------------
## This makefile can be used to execute common functions to interact with
## the source code, these functions ease local development and can also be
## used in CI/CD pipelines.
## ----------------------------------------------------------------------

godoc_version=v0.1.10
go_lint_version=v1.51.2
protoc_version=3.20.1
proto_c_gen_go_version=v1.28.0
proto_c_gen_go_grpc_version=v1.2
proto_linter_version=v0.50.5
swagger_version=v0.29.0

godoc_port=8081
swagger_port=8082

# REFERENCE: https://stackoverflow.com/questions/16931770/makefile4-missing-separator-stop
help: ## - Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

check-proto: ## - check protoc/proto-gen-go/protolint
	@which protolint || (go install github.com/yoheimuta/protolint/cmd/protolint@${proto_lint_version})
	@which protoc || echo protocv${protoc_version} not installed, install https://github.com/protocolbuffers/protobuf/releases/tag/v${protoc_version}
	@which protoc-gen-go || (go install google.golang.org/protobuf/cmd/protoc-gen-go@${proto_c_gen_go_version})
	@which protoc-gen-go-grpc || (go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${proto_c_gen_go_grpc_version})

lint-proto: check-proto ## - lint proto
	@protolint -fix ./data/pb/*.proto

clean-proto: ## - clean protos
	@rm ./data/pb/*.pb.go

build-proto: ## - build proto
	@protoc -I="./data/pb" --go_opt=paths=source_relative --go_out="./data/pb" --go-grpc_opt=paths=source_relative --go-grpc_out="./data/pb" ./data/pb/*.proto

check-lint: ## - validate/install golangci-lint installation
	which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@${proto_linter_version})

lint: check-lint ## - lint the source
	@golangci-lint run

lint-verbose: check-lint ## - lint the source with verbose output
	@golangci-lint run --verbose

# Reference: https://medium.com/@pedram.esmaeeli/generate-swagger-specification-from-go-source-code-648615f7b9d9
check-swagger: ## - validate/install swagger (v0.29.0)
	which swagger || (go install github.com/go-swagger/go-swagger/cmd/swagger@${swagger_version})

swagger: check-swagger ## - generate the swagger.json
	@swagger generate spec --work-dir=./internal/swagger --output ./tmp/swagger.json --scan-models

validate-swagger: swagger ## - validate the swagger.json
	@swagger validate ./tmp/swagger.json

serve-swagger: swagger ## - serve (web) the swagger.json
	@swagger serve -F=swagger ./tmp/swagger.json -p ${swagger_port} --no-open

check-godoc: ## - validate/install godoc
	which godoc || (go install golang.org/x/tools/cmd/godoc@${godoc_version})

serve-godoc: check-godoc ## - serve (web) the godocs
	@cd .. && godoc -http :${godoc_port}

test: run ## - test the source
	@go test -v -cover -parallel=1 --count=1 ./... -coverprofile ./tmp/go-bludgeon-healthcheck.out | tee ./tmp/go-bludgeon-healthcheck.log

build: ## - build the source (latest)
	@docker compose build --build-arg GIT_COMMIT=`git rev-parse HEAD` --build-arg GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
	@docker image prune -f

pull:
	@docker compose pull
	@docker image prune -f

dep:

run: ## - run the service and its dependencies (docker)
	@docker container rm -f healthcheck
	@docker image prune -f
	@docker compose --profile application up -d --wait

stop:
	@docker compose --profile application down

clean: stop ## - stop and clean docker resources
	@docker compose --profile application rm -f

install-client: ## - install the healthcheck client
	@go install ./cmd/healthcheck-client