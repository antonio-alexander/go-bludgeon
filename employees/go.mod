module github.com/antonio-alexander/go-bludgeon/employees

go 1.16

replace github.com/antonio-alexander/go-bludgeon/changes => ../changes

require (
	github.com/antonio-alexander/go-bludgeon/changes v0.0.0-20221225070951-9cd52c93b78a
	github.com/antonio-alexander/go-bludgeon/internal v1.2.2-0.20221225062935-a7f35fc350ec
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.0
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.0
)
