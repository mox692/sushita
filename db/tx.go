package db

import (
	"database/sql"
	"log"
)

/* トランザクション処理 */
func Transaction(txFunc func(*sql.Tx) error) error {
	tx, err := DbConnection.Begin()
	if err != nil {
		log.Println(err, "トランザクションの開始に失敗しました")
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Println(rollbackErr, "defer内のロールバックに失敗しました。")
			} else {
				log.Println("ロールバックしました。")
			}
			panic(p)
		} else if err != nil {
			log.Println(err)
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Println(rollbackErr, "defer内のロールバックに失敗しました。")
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
