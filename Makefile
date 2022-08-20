# Reference: https://medium.com/@pedram.esmaeeli/generate-swagger-specification-from-go-source-code-648615f7b9d9
check-swagger:
	@which swagger || (go install github.com/go-swagger/go-swagger/cmd/swagger@v0.29.0)

swagger: check-swagger
	@swagger generate spec --work-dir=./timers/service/rest/swagger --output ./tmp/swagger/swagger_timers.json --scan-models
	@swagger generate spec --work-dir=./employees/service/rest/swagger --output ./tmp/swagger/swagger_employees.json --scan-models
	@swagger generate spec --work-dir=./changes/internal/swagger --output ./tmp/swagger/swagger_changes.json --scan-models

serve-swagger: swagger
	@swagger mixin ./tmp/swagger/swagger_employees.json ./tmp/swagger/swagger_timers.json ./tmp/swagger/swagger_changes.json -o ./tmp/swagger/swagger.json
	@swagger serve -F=swagger ./tmp/swagger/swagger.json -p 8080 --no-open

export-swagger:
	@docker compose up -d
	@rm -rf ./tmp/swagger/html/*
	@mkdir -p ./tmp/swagger/html/api
	@docker exec -it swagger \
	@wget --mirror --adjust-extension --page-requisites \
	--no-host-directories --directory-prefix=/tmp/html \
	http://localhost:8080/
	@docker exec -it swagger wget --output-file /tmp/html/api/swagger_employees.json \
	http://localhost:8080/api/swagger_employees.json
	@docker exec -it swagger wget --output-file /tmp/html/api/swagger_timers.json \
	http://localhost:8080/api/swagger_timers.json
	@docker exec -it swagger wget --output-file /tmp/html/api/swagger_changes.json \
	http://localhost:8080/api/swagger_changes.json

check-godoc:
	@which godoc || (go install golang.org/x/tools/cmd/godoc@v0.1.10)

serve-godoc: check-godoc
	@godoc -http :8080

build:
	@docker compose build

run:
	@docker compose up

run-detach:
	@docker compose up -d

stop:
	@docker compose down

clean:
	@docker compose down
	@docker system prune -f