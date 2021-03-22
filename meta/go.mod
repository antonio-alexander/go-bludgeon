module github.com/antonio-alexander/go-bludgeon/meta

go 1.16

replace github.com/antonio-alexander/go-bludgeon/common => ../common

replace github.com/antonio-alexander/go-bludgeon/meta => ../meta

require (
	github.com/antonio-alexander/go-bludgeon/common v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.5.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
)
