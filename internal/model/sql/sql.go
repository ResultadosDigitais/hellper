package sql

import (
	"database/sql"
	"errors"
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
	//func (db *DB) Begin() (*Tx, error)
	//func (db *DB) Driver() driver.Driver
	//func (db *DB) Prepare(query string) (*Stmt, error)
	//func (db *DB) SetConnMaxLifetime(d time.Duration)
	//func (db *DB) SetMaxIdleConns(n int)
	//func (db *DB) SetMaxOpenConns(n int)
	//func (db *DB) Stats() DBStats
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
	return sql.Open(driver, dsn)
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
