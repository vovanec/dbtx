package dbtx

import (
	"context"
	"database/sql"
)

type preparedTransaction struct {
	stmt           string
	db             DB
	isolationLevel sql.IsolationLevel
	errorLog       ErrorLogFunc
}

func (t *preparedTransaction) Execute(f PreparedTxFunc) error {

	if f == nil {
		f = func(stmt *sql.Stmt) error {
			_, err := stmt.Exec()
			return err
		}
	}

	dbHandle, err := t.db.RWHandle()
	if err != nil {
		return err
	}

	tx, err := dbHandle.BeginTx(context.Background(), &sql.TxOptions{Isolation: t.isolationLevel})
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(t.stmt)
	defer func() {
		if err := stmt.Close(); err != nil {
			t.errorLog("error on statement close: %s", err)
		}
	}()

	if err := f(stmt); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			t.errorLog("rollback error: %s", rollbackErr)
		}
		return err
	}

	return tx.Commit()
}

func (t *preparedTransaction) ExecuteContext(ctx context.Context, f PreparedTxFunc) error {

	if f == nil {
		f = func(stmt *sql.Stmt) error {
			_, err := stmt.ExecContext(ctx)
			return err
		}
	}

	dbHandle, err := t.db.RWHandle()
	if err != nil {
		return err
	}

	tx, err := dbHandle.BeginTx(ctx, &sql.TxOptions{Isolation: t.isolationLevel})
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, t.stmt)
	defer func() {
		if err := stmt.Close(); err != nil {
			t.errorLog("error on statement close: %s", err)
		}
	}()

	if err := f(stmt); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			t.errorLog("rollback error: %s", rollbackErr)
		}
		return err
	}

	return tx.Commit()
}

func (t *preparedTransaction) WithIsolationLevel(l sql.IsolationLevel) PreparedTransaction {
	t.isolationLevel = l
	return t
}
