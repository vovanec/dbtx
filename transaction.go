package dbtx

import (
	"context"
	"database/sql"
)

type transaction struct {
	db             DB
	errorLog       ErrorLogFunc
	isolationLevel sql.IsolationLevel
}

func (t *transaction) Execute(f TxFunc) error {

	dbHandle, err := t.db.RWHandle()
	if err != nil {
		return err
	}

	tx, err := dbHandle.BeginTx(context.Background(), &sql.TxOptions{Isolation: t.isolationLevel})
	if err != nil {
		return err
	}

	if err := f(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			t.errorLog("rollback error: %s", rollbackErr)
		}
		return err
	}

	return tx.Commit()
}

func (t *transaction) ExecuteContext(ctx context.Context, f TxFunc) error {

	dbHandle, err := t.db.RWHandle()
	if err != nil {
		return err
	}

	tx, err := dbHandle.BeginTx(ctx, &sql.TxOptions{Isolation: t.isolationLevel})
	if err != nil {
		return err
	}

	if err := f(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			t.errorLog("rollback error: %s", rollbackErr)
		}
		return err
	}

	return tx.Commit()
}

func (t *transaction) WithIsolationLevel(l sql.IsolationLevel) Transaction {
	t.isolationLevel = l
	return t
}
