module github.com/antonio-alexander/go-bludgeon/meta

go 1.16

// replace github.com/antonio-alexander/go-bludgeon/internal => ../internal

replace github.com/antonio-alexander/go-bludgeon/data => ../data

require (
	github.com/antonio-alexander/go-bludgeon/data v1.0.0
	github.com/antonio-alexander/go-bludgeon/internal v1.0.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
)
