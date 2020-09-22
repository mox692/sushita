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

func CreateRankingtable(DbConnection *sql.DB, tx *sql.Tx) error {
	sqlCmd := `CREATE TABLE IF NOT EXISTS local_ranking(
		score INTEGER,
		created_at TEXT)`
	_, err := tx.Exec(sqlCmd)
	if err != nil {
		return fmt.Errorf("err : %w", err)
	}
	return nil
}

func InsertUserData(userID, userName string, tx *sql.Tx) error {
	// _, err := tx.Prepare("INSERT INTO user (id, user_name) values (?, ?);")
	// if err != nil {
	// 	return fmt.Errorf("db.Conn.Prepare err : %w", err)
	// }
	_, err := tx.Exec("INSERT INTO user (id, user_name) values (?, ?);", userID, userName)
	if err != nil {
		return fmt.Errorf("EXEC err : %w", err)
	}

	return nil
}
