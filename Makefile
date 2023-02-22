## ----------------------------------------------------------------------
## This makefile can be used to execute common functions to interact with
## the source code, these functions ease local development and can also be
## used in CI/CD pipelines.
## ----------------------------------------------------------------------

godoc_port = 9001 # defaults

check-swagger:
	@which swagger || (go install github.com/go-swagger/go-swagger/cmd/swagger@v0.29.0)

swagger: check-swagger
	@swagger generate spec --work-dir=./timers/internal/swagger --output ./tmp/swagger_timers.json --scan-models
	@swagger generate spec --work-dir=./employees/internal/swagger --output ./tmp/swagger_employees.json --scan-models
	@swagger generate spec --work-dir=./changes/internal/swagger --output ./tmp/swagger_changes.json --scan-models

serve-swagger: swagger
	@swagger mixin ./tmp/swagger_employees.json ./tmp/swagger_timers.json ./tmp/swagger_changes.json -o ./tmp/swagger.json
	@swagger serve -F=swagger ./tmp/swagger/swagger.json -p 8080 --no-open

check-godoc:
	@which godoc || (go install golang.org/x/tools/cmd/godoc@v0.1.10)

serve-godoc: check-godoc
	@godoc -http :${godoc_port}

build:
	@docker compose build --build-arg GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
	@docker image prune -f

run:
	@docker image prune -f
	@docker compose up -d --wait

stop:
	@docker compose down
