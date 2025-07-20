package database_connections

import (
	"database/sql"
)

type SQLAdapter interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Close() error
	Begin() (*sql.Tx, error)
}

type DefaultSQLAdapter struct {
	db *sql.DB
}

func NewSQLAdapter(db *sql.DB) SQLAdapter {
	return &DefaultSQLAdapter{db: db}
}

func (a *DefaultSQLAdapter) QueryRow(query string, args ...interface{}) *sql.Row {
	return a.db.QueryRow(query, args...)
}

func (a *DefaultSQLAdapter) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return a.db.Query(query, args...)
}

func (a *DefaultSQLAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return a.db.Exec(query, args...)
}

func (a *DefaultSQLAdapter) Close() error {
	return a.db.Close()
}

func (a *DefaultSQLAdapter) Begin() (*sql.Tx, error) {
	return a.db.Begin()
}

func SqlOpen(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(driverName, dataSourceName)
}
