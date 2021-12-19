module github.com/antonio-alexander/go-bludgeon/meta

go 1.16

replace github.com/antonio-alexander/go-bludgeon => ../

replace github.com/antonio-alexander/go-bludgeon/data => ../data

require (
	github.com/antonio-alexander/go-bludgeon/data v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.6.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
)
