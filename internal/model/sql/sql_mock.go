package sql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/stretchr/testify/mock"
)

type dynamicData map[string]interface{}

func (d dynamicData) Value() (driver.Value, error) {
	j, err := json.Marshal(d)
	return j, err
}

func (d *dynamicData) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("err_invalid_dbtype: != []byte")
	}

	err := json.Unmarshal(source, d)
	if err != nil {
		return err
	}
	return nil
}

func NewDBMock() *DBMock {
	return new(DBMock)
}

type DBMock struct {
	mock.Mock
}

func (mock *DBMock) QueryRow(sql string, params ...interface{}) Row {
	var (
		args   = mock.Called(sql, params)
		result = args.Get(0)
	)
	if result != nil {
		return result.(Row)
	}
	return nil
}

func (mock *DBMock) Query(sql string, params ...interface{}) (Rows, error) {
	var (
		args   = mock.Called(sql, params)
		result = args.Get(0)
	)
	if result != nil {
		return result.(Rows), args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *DBMock) Exec(query string, params ...interface{}) (Result, error) {
	var (
		args   = mock.Called(query, params)
		result = args.Get(0)
	)
	if result != nil {
		return result.(Result), args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *DBMock) Ping() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *DBMock) Close() error {
	args := mock.Called()
	return args.Error(0)
}

func NewRowMock() *RowMock {
	return new(RowMock)
}

type RowMock struct {
	mock.Mock
}

func (mock *RowMock) Scan(dest ...interface{}) error {
	args := mock.Called(dest)
	return args.Error(0)
}

func NewRowsMock() *RowsMock {
	return new(RowsMock)
}

type RowsMock struct {
	mock.Mock
}

func (mock *RowsMock) Next() bool {
	args := mock.Called()
	return args.Bool(0)
}

func (mock *RowsMock) Scan(dest ...interface{}) error {
	args := mock.Called(dest)
	return args.Error(0)
}

func NewResultMock() *ResultMock {
	return new(ResultMock)
}

type ResultMock struct {
	mock.Mock
}

func (mock *ResultMock) LastInsertId() (int64, error) {
	args := mock.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (mock *ResultMock) RowsAffected() (int64, error) {
	args := mock.Called()
	return args.Get(0).(int64), args.Error(1)
}
