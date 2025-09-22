package orm

import (
	"boilerplate_backend_api_gin/core/database/database_connections"
	"database/sql"
)

// Transaction maneja transacciones de base de datos
type Transaction struct {
	tx *sql.Tx
}

// TxAdapter es un adaptador para usar una transacción SQL
type TxAdapter struct {
	tx *sql.Tx
}

func (a *TxAdapter) QueryRow(query string, args ...interface{}) *sql.Row {
	return a.tx.QueryRow(query, args...)
}

func (a *TxAdapter) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return a.tx.Query(query, args...)
}

func (a *TxAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return a.tx.Exec(query, args...)
}

func (a *TxAdapter) Close() error {
	// La transacción no se cierra directamente, se maneja con Commit/Rollback
	return nil
}

func (a *TxAdapter) Begin() (*sql.Tx, error) {
	// No se puede iniciar una transacción dentro de otra
	return nil, nil
}

// NewTransaction inicia una nueva transacción
func NewTransaction() (*Transaction, error) {
	db := database_connections.DatabaseConnectSQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	return &Transaction{tx: tx}, nil
}

// Commit confirma la transacción
func (t *Transaction) Commit() error {
	return t.tx.Commit()
}

// Rollback revierte la transacción
func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}

// GetAdapter devuelve un adaptador SQL para usar con la transacción
func (t *Transaction) GetAdapter() database_connections.SQLAdapter {
	return &TxAdapter{tx: t.tx}
}

// WithTransaction ejecuta una función dentro de una transacción
func WithTransaction(fn func(*Transaction) error) error {
	tx, err := NewTransaction()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// Asegurarse de hacer rollback en caso de pánico
			_ = tx.Rollback()
			panic(p) // Re-panic después del rollback
		}
	}()

	err = fn(tx)

	if err != nil {
		// Rollback en caso de error
		_ = tx.Rollback()
		return err
	}

	// Commit si todo salió bien
	return tx.Commit()
}
