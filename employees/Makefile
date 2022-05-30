check-lint:
	which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2)

lint: check-lint
	golangci-lint run

lint-verbose: check-lint
	golangci-lint run --verbose

# Reference: https://medium.com/@pedram.esmaeeli/generate-swagger-specification-from-go-source-code-648615f7b9d9
check-swagger:
	which swagger || (go install github.com/go-swagger/go-swagger/cmd/swagger@v0.29.0)

swagger: check-swagger
	swagger generate spec -o ./cmd/service/swagger.yaml --scan-models

serve-swagger: swagger
	swagger serve -F=swagger ./cmd/service/swagger.yaml -p 8080 --no-open

check-godoc:
	which godoc || (go install golang.org/x/tools/cmd/godoc@v0.1.10)

serve-godoc: check-godoc
	cd .. && godoc -http :8080

test:
	 go test -cover --count=1 ./...

test-verbose:
	 go test -cover -v --count=1 ./...

build:
	./build.sh latest

run: build
	docker compose up