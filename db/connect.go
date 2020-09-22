package db

import (
	"database/sql"
	"fmt"
	"os"
)

var DbConnection *sql.DB

func DBinit() error {
	workDir, _ := os.Getwd()
	dbPath := workDir + "/db.sql"
	var err error
	DbConnection, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("err : %w", err)
	}
	return nil
}
