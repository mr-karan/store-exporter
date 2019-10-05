package store

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Imports mysql driver
	_ "github.com/lib/pq"              // Imports postgresql driver
	_ "github.com/mattn/go-sqlite3"    // Imports sqlite driver
)

// Store represents additional options for the external store
type Store struct {
	MaxOpenConnections int    `koanf:"max_open_connections"`
	MaxIdleConnections int    `koanf:"max_idle_connections"`
	DB                 string `koanf:"db"`
	DSN                string `koanf:"dsn"`
	QueryFile          string `koanf:"query"`
}

// Manager represents the set of methods used to interact with the db.
type Manager interface {
	FetchResults(string) (map[string]interface{}, error)
}

// DBConnOpts represents additonal parameters to create a DB Client

// NewManager instantiates an object of Manager based on the params
func NewManager(store Store) (Manager, error) {
	switch dbType := store.DB; dbType {
	case "postgres", "mysql":
		return NewDBClient(store.DB, store.DSN, &DBClientOpts{
			QueryFile:    store.QueryFile,
			MaxIdleConns: store.MaxIdleConnections,
			MaxOpenConns: store.MaxOpenConnections,
		})
	case "sqlite3":
		return NewDBClient(store.DB, store.DSN, &DBClientOpts{
			QueryFile:    store.QueryFile,
			MaxIdleConns: store.MaxIdleConnections,
			MaxOpenConns: store.MaxOpenConnections,
		})
	// TODO case "redis":
	default:
		return nil, fmt.Errorf("Error fetching results: Unknown db type")
	}
}
