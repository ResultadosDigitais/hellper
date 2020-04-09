package sql

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

type testDB struct {
	name string
	db   *sql.DB
	err  error
}

func (scenario *testDB) setup(t *testing.T) {
	if scenario.err == nil {
		db, mock, err := sqlmock.New()
		require.NotNil(t, db, "db instance")
		require.NotNil(t, mock, "mock db instance")
		require.Nil(t, err, "sqlmock error")

		if scenario.err == nil {
			mock.ExpectClose()
		}

		scenario.db = db
	}
}

func (scenario *testDB) tearDown(t *testing.T) {
	if scenario.db != nil {
		scenario.db.Close()
	}
}

func TestDB(test *testing.T) {
	scenarios := []testDB{
		{
			name: "Creates a new client",
		},
		{
			name: "Returns error because db is blank",
			err:  ErrBlankDB,
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				db, err := NewDB(scenario.db)
				require.Equal(t, scenario.err, err, "newDB error")
				if scenario.err == nil {
					require.NotNil(t, db, "db instance")
					require.Nil(t, db.Ping(), "ping error")
					require.Nil(t, db.Close(), "close error")
				} else {
					require.Nil(t, db, "db invalid instance")
				}
			},
		)
	}
}

type testQuery struct {
	name      string
	db        *sql.DB
	sqlMock   sqlmock.Sqlmock
	query     string
	arguments []interface{}
	columns   []string
	rows      [][]interface{}
	err       error
	rowsErr   error
}

func (scenario *testQuery) setup(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NotNil(t, db, "db instance")
	require.NotNil(t, mock, "mock db instance")
	require.Nil(t, err, "sqlmock error")

	if scenario.err == nil {
		mockRows := sqlmock.NewRows(scenario.columns)
		if scenario.rowsErr == nil {
			for _, row := range scenario.rows {
				columns := make([]driver.Value, len(row))
				for index, column := range row {
					columns[index] = column
				}
				mockRows.AddRow(columns...)
			}
		} else {
			mockRows.RowError(1, scenario.rowsErr)
		}
		mock.ExpectQuery(scenario.query).WillReturnRows(mockRows)
	} else {
		mock.ExpectQuery(scenario.query).WillReturnError(scenario.err)
	}

	mock.ExpectClose()

	scenario.db = db
	scenario.sqlMock = mock
}

func (scenario *testQuery) tearDown(t *testing.T) {
	if scenario.db != nil {
		scenario.db.Close()
	}
}

func TestQuery(test *testing.T) {
	scenarios := []testQuery{
		{
			name:    "Executes query successfully",
			query:   "select id, text from mock",
			columns: []string{"id", "text"},
			rows: [][]interface{}{
				{int64(1), "mock one"},
				{int64(2), "mock two"},
			},
		},
		{
			name:    "When query returns error",
			query:   "select",
			columns: []string{"id", "text"},
			err:     errors.New("err_mockquery"),
		},
		{
			name:    "Returns error when try to read results",
			query:   "select",
			columns: []string{"id", "text"},
			rowsErr: errors.New("err_mockrows"),
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				db, err := NewDB(scenario.db)
				require.Nil(t, err, "newDB error")
				require.NotNil(t, db, "db instance")
				rows, err := db.Query(scenario.query, scenario.arguments...)
				require.Equal(t, scenario.err, err, "query error")
				if scenario.err == nil {
					for index := 0; rows.Next(); index++ {
						var (
							values   = make([]interface{}, len(scenario.columns))
							pointers = make([]interface{}, len(scenario.columns))
						)
						for index := range scenario.columns {
							pointers[index] = &values[index]
						}
						err := rows.Scan(pointers...)
						require.Equal(t, scenario.rowsErr, err, "rows error")
						if scenario.rowsErr == nil {
							row := scenario.rows[index]
							require.Equal(t, row, values, "row error")
						}
					}
				} else {
					require.Nil(t, rows, "rows invalid instance")
				}
				require.Nil(t, db.Close(), "close error")
				require.Nil(t, scenario.sqlMock.ExpectationsWereMet(), "sqlmock invalid expectations")
			},
		)
	}
}

type testQueryRow struct {
	name      string
	db        *sql.DB
	sqlMock   sqlmock.Sqlmock
	query     string
	arguments []interface{}
	columns   []string
	row       []interface{}
	err       error
}

func (scenario *testQueryRow) setup(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NotNil(t, db, "db instance")
	require.NotNil(t, mock, "mock db instance")
	require.Nil(t, err, "sqlmock error")

	mockRows := sqlmock.NewRows(scenario.columns)
	if scenario.err == nil {
		columns := make([]driver.Value, len(scenario.row))
		for index, column := range scenario.row {
			columns[index] = column
		}
		mockRows.AddRow(columns...)
		mock.ExpectQuery(scenario.query).WillReturnRows(mockRows)
	} else {
		mock.ExpectQuery(scenario.query).WillReturnError(scenario.err)
	}

	mock.ExpectClose()

	scenario.db = db
	scenario.sqlMock = mock
}

func (scenario *testQueryRow) tearDown(t *testing.T) {
	if scenario.db != nil {
		scenario.db.Close()
	}
}

func TestQueryRow(test *testing.T) {
	scenarios := []testQueryRow{
		{
			name:    "Executes query row successfully",
			query:   "select id, text from mock",
			columns: []string{"id", "text"},
			row: []interface{}{
				int64(1), "mock one",
			},
		},
		{
			name:    "Returns error when try to read results",
			query:   "select",
			columns: []string{"id", "text"},
			err:     errors.New("err_mockrows"),
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				db, err := NewDB(scenario.db)
				require.Nil(t, err, "newDB error")
				require.NotNil(t, db, "db instance")
				row := db.QueryRow(scenario.query, scenario.arguments...)
				require.NotNil(t, row, "row invalid instance")
				var (
					values   = make([]interface{}, len(scenario.columns))
					pointers = make([]interface{}, len(scenario.columns))
				)
				for index := range scenario.columns {
					pointers[index] = &values[index]
				}
				err = row.Scan(pointers...)
				require.Equal(t, scenario.err, err, "scan error")
				if scenario.err == nil {
					require.Equal(t, scenario.row, values, "row columns invalid instance")
				}
				require.Nil(t, db.Close(), "close error")
				require.Nil(t, scenario.sqlMock.ExpectationsWereMet(), "sqlmock invalid expectations")
			},
		)
	}
}

type testExec struct {
	name         string
	db           *sql.DB
	sqlMock      sqlmock.Sqlmock
	query        string
	arguments    []interface{}
	lastID       int64
	rowsAffected int64
	err          error
}

func (scenario *testExec) setup(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NotNil(t, db, "db instance")
	require.NotNil(t, mock, "mock db instance")
	require.Nil(t, err, "sqlmock error")

	if scenario.err == nil {
		mockResult := sqlmock.NewResult(scenario.lastID, scenario.rowsAffected)
		mock.ExpectExec(scenario.query).WillReturnResult(mockResult)
	} else {
		mock.ExpectExec(scenario.query).WillReturnError(scenario.err)
	}

	mock.ExpectClose()

	scenario.db = db
	scenario.sqlMock = mock
}

func (scenario *testExec) tearDown(t *testing.T) {
	if scenario.db != nil {
		scenario.db.Close()
	}
}

func TestExec(test *testing.T) {
	scenarios := []testExec{
		{
			name:  "Executes command successfully",
			query: "insert into mock",
			arguments: []interface{}{
				1, "mock one",
			},
			lastID:       1,
			rowsAffected: 1,
		},
		{
			name:  "Returns error when try to execute command",
			query: "insert into mock",
			err:   errors.New("err_mockrows"),
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				db, err := NewDB(scenario.db)
				require.Nil(t, err, "newDB error")
				require.NotNil(t, db, "db instance")
				result, err := db.Exec(scenario.query, scenario.arguments...)
				require.Equal(t, scenario.err, err, "exec error")
				if scenario.err == nil {
					require.NotNil(t, result, "result invalid instance")

					lastID, err := result.LastInsertId()
					require.Nil(t, err, "lastinsertid error")
					require.Equal(t, scenario.lastID, lastID, "lastinsertid invalid instance")

					rowsAffected, err := result.RowsAffected()
					require.Nil(t, err, "rowsaffected error")
					require.Equal(t, scenario.rowsAffected, rowsAffected, "rowsaffected invalid instance")
				} else {
					require.Nil(t, result, "result invalid instance")
				}
				require.Nil(t, db.Close(), "close error")
				require.Nil(t, scenario.sqlMock.ExpectationsWereMet(), "sqlmock invalid expectations")
			},
		)
	}
}
