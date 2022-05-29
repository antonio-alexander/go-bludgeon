# Reference: https://medium.com/@pedram.esmaeeli/generate-swagger-specification-from-go-source-code-648615f7b9d9
check-swagger:
	which swagger || (go install github.com/go-swagger/go-swagger/cmd/swagger@v0.29.0)

swagger: check-swagger
	swagger generate spec -o ./cmd/service/swagger.yaml --scan-models

serve-swagger: check-swagger
	swagger serve -F=swagger ./cmd/service/swagger.yaml

check-godoc:
	which godoc || (go install golang.org/x/tools/cmd/godoc@v0.1.10)

serve-godoc: check-godoc
	cd .. && godoc -http :8080