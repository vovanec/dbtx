package dbtx

import (
	"database/sql"
	"fmt"
	"github.com/vovanec/xsync"
)

type db struct {
	conf       Config
	reader     *sql.DB
	writer     *sql.DB
	readerOnce xsync.Once
	writerOnce xsync.Once
	errorLog   ErrorLogFunc
}

func (d *db) RWHandle() (*sql.DB, error) {

	var (
		err    error
		handle *sql.DB
	)

	d.writerOnce.Do(func() error {
		handle, err = d.getDBHandle(true)
		if err != nil {
			return err
		}
		d.writer = handle

		return nil
	})

	return d.writer, err
}

func (d *db) ROHandle() (*sql.DB, error) {

	var (
		err    error
		handle *sql.DB
	)

	d.readerOnce.Do(func() error {
		handle, err = d.getDBHandle(false)
		if err != nil {
			return err
		}
		d.reader = handle

		return nil
	})

	return d.reader, err
}

func (d *db) PreparedTransaction(stmt string) PreparedTransaction {

	return &preparedTransaction{
		db:   d,
		stmt: stmt,
		errorLog:  d.errorLog,
	}
}

func (d *db) Transaction() Transaction {

	return &transaction{
		db:  d,
		errorLog: d.errorLog,
	}
}

func (d *db) getDBHandle(isWriter bool) (*sql.DB, error) {

	var connStr string
	if isWriter {
		connStr = d.conf.GetRWDataSourceName()
	} else {
		connStr = d.conf.GetRODataSourceName()
	}

	db, err := sql.Open(d.conf.DriverName(), connStr)
	if err != nil {
		return nil, fmt.Errorf("could not get database connection: %s", err)
	}

	if connPoolConf, ok := d.conf.(ConnPoolConfig); ok {
		db.SetMaxOpenConns(connPoolConf.MaxOpenConnections())
		db.SetMaxIdleConns(connPoolConf.MaxIdleConnections())
	}


	return db, nil
}

func NewDB(conf Config) DB {

	var errLogFunc ErrorLogFunc

	if errLogConf, ok := conf.(ErrorLogConfig); ok {
		errLogFunc = errLogConf.ErrorLog()
	} else {
		errLogFunc = func (format string, args ...interface{}) {}
	}

	return &db{
		conf: conf,
		errorLog:  errLogFunc,
	}
}

