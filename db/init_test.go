package db_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mox692/sushita/cmd"
)

// useridとusernameを適当に与えて、
// mockのexpect通りの結果が得られるか。
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

	err = cmd.SetupDB("testID", "testName", db)

}
