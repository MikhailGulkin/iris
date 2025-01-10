module chat

go 1.23.4

require (
	github.com/gocql/gocql v0.0.0-20211015133455-b225f9b53fa1
	github.com/google/uuid v1.6.0
	github.com/scylladb/gocqlx/v3 v3.0.1
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/scylladb/go-reflectx v1.0.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

replace github.com/gocql/gocql => github.com/scylladb/gocql v1.14.4
