package db

import (
	"database/sql"
	"fmt"
)

func CreateUsertable(DbConnection *sql.DB) error {
	sqlCmd := `CREATE TABLE IF NOT EXISTS user(
		id TEXT PRIMARY KEY,
		user_name TEXT)`
	_, err := DbConnection.Exec(sqlCmd)
	if err != nil {
		return fmt.Errorf("err : %w", err)
	}
	return nil
}

func CreateRankingtable(DbConnection *sql.DB) error {
	sqlCmd := `CREATE TABLE IF NOT EXISTS local_ranking(
		score INTEGER,
		created_at TEXT)`
	_, err := DbConnection.Exec(sqlCmd)
	if err != nil {
		return fmt.Errorf("err : %w", err)
	}
	return nil
}
