package store

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// handleTransactions will iterate through a variadic argument of transactions
// If there is at least one transaction, use that transaction for the query
// If there are no transactions, wrap the query in a transaction that begins, rollbacks, and commits
func handleTransactions(conn *sqlx.DB, fn func(tx *sql.Tx) error, txes ...*sql.Tx) error {
	var err error
	if len(txes) > 0 {
		tx := txes[0]
		err = fn(tx)
		if err != nil {
			return err
		}
		return nil
	}

	err = transact(conn, func(tx *sql.Tx) error {
		return fn(tx)
	})
	return err
}

// transact will wrap a query with begin, rollback and commit funcs
func transact(conn *sqlx.DB, fn func(tx *sql.Tx) error) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	err = fn(tx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	tx.Rollback()
	return err
}
