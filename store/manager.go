package store

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Imports mysql driver
	_ "github.com/lib/pq"              // Imports postgresql driver
)

// Manager represents the set of methods used to interact with the db.
type Manager interface {
	FetchResults(string) (map[string]interface{}, error)
}

// DBConnOpts represents additonal parameters to create a DB Client
type DBConnOpts struct {
	QueryFilePath string
}

// NewManager instantiates an object of Manager based on the params
func NewManager(db string, dsn string, opts *DBConnOpts) (Manager, error) {
	switch dbType := db; dbType {
	case "postgres":
		return NewDBClient(db, dsn, opts.QueryFilePath)
	case "mysql":
		return NewDBClient(db, dsn, opts.QueryFilePath)
	// TODO case "redis":
	default:
		return nil, fmt.Errorf("Error fetching results: Unknown db type")
	}
}
