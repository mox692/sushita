package db

import (
	"database/sql"
	"fmt"
	"os/user"

	_ "github.com/mattn/go-sqlite3"
)

var DbConnection *sql.DB

func DBinit() error {

	user, err := user.Current()
	if err != nil {
		return fmt.Errorf("err : %w", err)
	}
	dbPath := user.HomeDir + "/db.sql"
	DbConnection, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("err : %w", err)
	}
	return nil
}
