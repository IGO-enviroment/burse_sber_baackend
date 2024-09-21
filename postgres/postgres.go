package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type PostgresDbConnector struct {
	connectionString string
	db               *sql.DB
}

func NewPostgresConnector(connectionString string) *PostgresDbConnector {
	return &PostgresDbConnector{connectionString: connectionString}
}

func (c *PostgresDbConnector) Open() (*sql.DB, error) {
	db, err := sql.Open("postgres", c.connectionString)
	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, fmt.Errorf("Cannot connect to postgres, connection string: %s\n", c.connectionString)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	c.db = db
	fmt.Printf("Connected to Postgres: %s\n", c.connectionString)
	return c.db, nil
}

func (c *PostgresDbConnector) Close() error {
	if err := c.db.Close(); err != nil {
		return err
	}
	fmt.Printf("Disconnected from Postgres\n")
	return nil
}
