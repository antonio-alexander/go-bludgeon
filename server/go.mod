module github.com/antonio-alexander/go-bludgeon/server

go 1.16

replace github.com/antonio-alexander/go-bludgeon/common => ../common

replace github.com/antonio-alexander/go-bludgeon/meta => ../meta

replace github.com/antonio-alexander/go-bludgeon/logic => ../logic

replace github.com/antonio-alexander/go-bludgeon/internal/rest => ../internal/rest

replace github.com/antonio-alexander/go-bludgeon/internal/logger => ../internal/logger

require (
	github.com/antonio-alexander/go-bludgeon/common v0.0.0-00010101000000-000000000000
	github.com/antonio-alexander/go-bludgeon/internal/logger v0.0.0-00010101000000-000000000000
	github.com/antonio-alexander/go-bludgeon/internal/rest v0.0.0-00010101000000-000000000000
	github.com/antonio-alexander/go-bludgeon/logic v0.0.0-00010101000000-000000000000
	github.com/antonio-alexander/go-bludgeon/meta v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
)
