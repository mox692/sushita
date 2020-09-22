package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// *sql.DBを入れた時に、tableが作成されるか
// func TestCreateUsertable(t *testing.T) {

// }

// useridとusernameを適当に与えて、
// mockのexpect通りの結果が得られるか。
func TestInsertUserData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO user").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}
