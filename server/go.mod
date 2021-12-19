module github.com/antonio-alexander/go-bludgeon/server

go 1.16

// replace github.com/antonio-alexander/go-bludgeon => ../

// replace github.com/antonio-alexander/go-bludgeon/data => ../data

require (
	github.com/antonio-alexander/go-bludgeon/data v1.0.0
	github.com/antonio-alexander/go-bludgeon/internal v1.0.0
	github.com/antonio-alexander/go-bludgeon/meta v1.0.0
	github.com/gorilla/mux v1.8.0
	github.com/pkg/errors v0.9.1
)
