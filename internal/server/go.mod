module github.com/antonio-alexander/go-bludgeon/internal/server

go 1.16

replace github.com/antonio-alexander/go-bludgeon/internal/common => ../common

replace github.com/antonio-alexander/go-bludgeon/internal/meta => ../meta

require (
	github.com/antonio-alexander/go-bludgeon/internal/common v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/pkg/errors v0.9.1
)
