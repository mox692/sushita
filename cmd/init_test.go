package cmd

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// useridとusernameを適当に与えて、
// SetupDBが正しく動作するか
func TestSetupDB(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS user").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS local_ranking").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO user").WithArgs("testID", "testName").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err = SetupDB("testID", "testName", db); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
