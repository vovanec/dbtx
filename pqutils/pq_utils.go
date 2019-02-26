package pqutils

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/vovanec/dbtx"
)

const (
	uniqueViolationErr       = "23505"
	unknownErr               = "unknown"
	savepointQuery           = "SAVEPOINT before_insert"
	rollbackToSavepointQuery = "ROLLBACK TO before_insert"
)

func TxUpsert(ctx context.Context, tx *sql.Tx, insertQuery, updateQuery string, args ...interface{}) error {

	var (
		err error
	)

	if _, err = tx.ExecContext(ctx, savepointQuery); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, insertQuery, args...); err != nil {

		if IsUniqueViolationErr(err) {
			if _, err = tx.ExecContext(ctx, rollbackToSavepointQuery); err != nil {
				return err
			}
			if _, err = tx.ExecContext(ctx, updateQuery, args...); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func DbUpsert(ctx context.Context, db dbtx.DB, insertQuery, updateQuery string, args ...interface{}) error {

	return db.Transaction().ExecuteContext(ctx,
		func(tx *sql.Tx) error {
			return TxUpsert(ctx, tx, insertQuery, updateQuery, args...)
		},
	)
}

func ErrorCode(err error) pq.ErrorCode {

	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code
	}

	return unknownErr
}

func IsUniqueViolationErr(err error) bool {

	return ErrorCode(err) == uniqueViolationErr
}
