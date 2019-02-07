package dbtx

import (
	"context"
	"database/sql"
)

type ErrorLogFunc func(format string, args ...interface{})

type DB interface {
	ROHandle() (*sql.DB, error)
	RWHandle() (*sql.DB, error)
	PreparedTransaction(statement string) PreparedTransaction
	Transaction() Transaction
}

type PreparedTxFunc func(stmt *sql.Stmt) error

type PreparedTransaction interface {
	ExecuteContext(ctx context.Context, f PreparedTxFunc) error
	Execute(f PreparedTxFunc) error
	WithIsolationLevel(level sql.IsolationLevel) PreparedTransaction
}

type TxFunc func(tx *sql.Tx) error

type Transaction interface {
	ExecuteContext(ctx context.Context, f TxFunc) error
	Execute(f TxFunc) error
	WithIsolationLevel(level sql.IsolationLevel) Transaction
}

type Config interface {
	DriverName() string
	GetRODataSourceName() string
	GetRWDataSourceName() string
}

type ConnPoolConfig interface {
	MaxOpenConnections() int
	MaxIdleConnections() int
}

type ErrorLogConfig interface {
	ErrorLog() ErrorLogFunc
}
