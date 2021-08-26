package storage

import (
	"database/sql"
	"fmt"
	// postgres lib
	_ "github.com/lib/pq"
	"lahaus/config"
)

// PostgreSQLManager ...
type PostgreSQLManager struct {
	Conn             *sql.DB
	ConnectionString string
}

// NewPostgreSQLManager ...
func NewPostgreSQLManager(config *config.Database) (*PostgreSQLManager, error) {
	manager := &PostgreSQLManager{}
	manager.ConnectionString = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.DatabaseName)
	conn, err := sql.Open("postgres", manager.ConnectionString)
	if err != nil {
		return nil, err
	}
	manager.Conn = conn
	if err := manager.Conn.Ping(); err != nil {
		return nil, err
	}
	return manager, nil
}
