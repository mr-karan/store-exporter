package store

import (
	"context"
	"database/sql"
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

// NewDBClient initializes a connection object with the databse.
func NewDBClient(db string, dsn string, queryFile string) (Manager, error) {
	conn, err := sqlx.Connect(db, dsn)
	if err != nil {
		return nil, err
	}
	// Connection Pool grows unbounded, have some sane defaults.
	conn.SetMaxIdleConns(3)
	conn.SetMaxOpenConns(5)
	// Load queries
	if queryFile == "" {
		return nil, fmt.Errorf("error initialising DB Manager: Path to query file not provided")
	}
	queries := goyesql.MustParseFile(queryFile)

	return &DBClient{
		Conn:    conn,
		Queries: queries,
	}, nil
}

// FetchResults executes the query and parses the result
func (client *DBClient) FetchResults(query string) (map[string]interface{}, error) {
	tx, err := client.Conn.BeginTxx(context.Background(), &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	q, ok := client.Queries[query]
	if !ok {
		return nil, fmt.Errorf("No query mapped to: %s", query)
	}
	row := tx.QueryRowx(q.Query)
	results := make(map[string]interface{})
	err = row.MapScan(results)
	if err != nil {
		return nil, err
	}
	return results, err
}
