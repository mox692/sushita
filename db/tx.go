package db

import (
	"database/sql"
	"log"
)

/* トランザクション処理 */
func Transaction(txFunc func(*sql.Tx) error, dbConn *sql.DB) error {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err, "fail dbConn.Begin")
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Println(rollbackErr, "fail tx.Rollback")
			} else {
				log.Println("ロールバックしました。")
			}
			panic(p)
		} else if err != nil {
			log.Println(err)
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Println(rollbackErr, "fail tx.Rollback")
			} else {
				log.Println("ロールバックしました。")
			}
		} else {
			err = tx.Commit()
			if err != nil {
				log.Println(err, "commitに失敗しました。")
				rollbackErr := tx.Rollback()
				if rollbackErr != nil {
					log.Println(rollbackErr, "rollbackに失敗しました")
				} else {
					log.Println("ロールバックしました。")
				}
			}
		}
	}()
	err = txFunc(tx)
	return err
}
