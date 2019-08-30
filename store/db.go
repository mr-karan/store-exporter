package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/goyesql"
)

// DBClient represents the structure to hold DB Client object required to create DB session and
// query and fetch results.
type DBClient struct {
	Conn    *sqlx.DB
	Queries goyesql.Queries
}

// DBClientOpts represents additional options to use for DB Client
type DBClientOpts struct {
	QueryFile    string
	MaxIdleConns int
	MaxOpenConns int
}

// NewDBClient initializes a connection object with the databse.
func NewDBClient(db string, dsn string, opts *DBClientOpts) (Manager, error) {
	conn, err := sqlx.Connect(db, dsn)
	if err != nil {
		return nil, err
	}
	// Connection Pool grows unbounded, have some sane defaults.
	conn.SetMaxIdleConns(opts.MaxIdleConns)
	conn.SetMaxOpenConns(opts.MaxOpenConns)
	// Load queries
	if opts.QueryFile == "" {
		return nil, fmt.Errorf("error initialising DB Manager: Path to query file not provided")
	}
	queries := goyesql.MustParseFile(opts.QueryFile)

	return &DBClient{
		Conn:    conn,
		Queries: queries,
	}, nil
}

// FetchResults executes the query and parses the result
func (client *DBClient) FetchResults(query string) (map[string]interface{}, error) {
	q, ok := client.Queries[query]
	if !ok {
		return nil, fmt.Errorf("No query mapped to: %s", query)
	}
	row := client.Conn.QueryRowx(q.Query)
	results := make(map[string]interface{})
	err := row.MapScan(results) // connection is closed automatically here. Read more: https://jmoiron.github.io/sqlx/#queryrow
	if err != nil {
		return nil, err
	}
	return results, err
}
