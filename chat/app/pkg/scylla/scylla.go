package scylla

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
)

type Config struct {
	Hosts []string
}

func NewScyllaClient(config Config) (*gocqlx.Session, error) {
	cluster := gocql.NewCluster(config.Hosts...)
	// Wrap session on creation, gocqlx session embeds gocql.Session pointer.
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return nil, err
	}
	return &session, nil
}
