package sql

import (
	"database/sql"
	"errors"
	"hellper/internal/config"
)

var (
	ErrBlankDB = errors.New("err_blankdb")
)

type DB interface {
	Query(string, ...interface{}) (Rows, error)
	QueryRow(string, ...interface{}) Row
	Exec(string, ...interface{}) (Result, error)
	Ping() error
	Close() error
}

type Row interface {
	Scan(...interface{}) error
}

type Rows interface {
	Row
	Next() bool
	Close() error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type db struct {
	*sql.DB
}

func (db *db) Query(sql string, arguments ...interface{}) (Rows, error) {
	rows, err := db.DB.Query(sql, arguments...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (db *db) QueryRow(sql string, arguments ...interface{}) Row {
	return db.DB.QueryRow(sql, arguments...)
}

func (db *db) Exec(sql string, arguments ...interface{}) (Result, error) {
	return db.DB.Exec(sql, arguments...)
}

func newSQLDB(driver, dsn string) (*sql.DB, error) {

	var (
		connectionSSLMode = config.Env.ConnectionSSLMode
	)

	psqlInfo := dsn
	if connectionSSLMode {
		sslmode := "&sslmode=verify-ca"
		sslrootcert := "&sslrootcert=/var/hellper/server-ca.pem"
		sslkey := "&sslkey=/var/hellper/client-key.pem"
		sslcert := "&sslcert=/var/hellper/client-cert.pem"
		psqlInfo = sslmode + sslrootcert + sslkey + sslcert
	}

	return sql.Open(driver, psqlInfo)
}

func NewDB(sqlDB *sql.DB) (DB, error) {
	if sqlDB == nil {
		return nil, ErrBlankDB
	}
	return &db{DB: sqlDB}, nil
}

func NewDBWithDSN(driver, dsn string) DB {
	sqlDB, err := newSQLDB(driver, dsn)
	if err != nil {
		panic(err)
	}
	db, err := NewDB(sqlDB)
	if err != nil {
		panic(err)
	}
	return db
}
