package db

import (
	"database/sql"
	"fmt"
)

type LocalRanking struct {
	Score     int
	CreatedAt string
}

type User struct {
	Id       int
	UserName string
}

func CreateUsertable(tx *sql.Tx) error {
	sqlCmd := `CREATE TABLE IF NOT EXISTS user(
		id TEXT PRIMARY KEY,
		user_name TEXT)`
	_, err := tx.Exec(sqlCmd)
	if err != nil {
		return fmt.Errorf("err : %w", err)
	}
	return nil
}

func CreateRankingtable(tx *sql.Tx) error {
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
	_, err := tx.Exec("INSERT INTO user (id, user_name) values (?, ?);", userID, userName)
	if err != nil {
		return fmt.Errorf("EXEC err : %w", err)
	}
	return nil
}

func SelectUser() (*User, error) {
	user := User{}
	row := DbConnection.QueryRow("select * from user limit 1;")
	err := row.Scan(&user.Id, &user.UserName)
	return &user, err
}
