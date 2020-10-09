package cmd

import (
	"database/sql"
	"fmt"

	"github.com/mox692/sushita/db"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "First command to start `sushita`",
	Long: `You have to run this command at first, 'sushita init'.\n
			This command will create the local-strage file db.sql\n
			(which store your score, userID for example) in your homedir.\n
			If you don't install sqlite3, this command will return error.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var userName string

		fmt.Printf("Enter your name!\n")
		fmt.Printf("=>")
		fmt.Scanf("%s", &userName)
		fmt.Printf("hello %s!!\n", userName)

		// UUIDでユーザIDを生成する
		userid, err := uuid.NewRandom()
		userID := userid.String()
		if err != nil {
			return fmt.Errorf("err : %w", err)
		}

		err = db.DBinit()
		if err != nil {
			fmt.Errorf("DbConnection.Exec : %w", err)
		}

		err = SetupDB(userID, userName, db.DbConnection)
		if err != nil {
			return fmt.Errorf("err : %w", err)
		}

		cmd.Printf("init sushita!🎉\n")
		return nil
	},
}

func SetupDB(userID, userName string, dbConn *sql.DB) error {

	err := db.Transaction(func(tx *sql.Tx) error {

		err := db.CreateUsertable(tx)
		if err != nil {
			return fmt.Errorf("err : %w", err)
		}

		err = db.CreateRankingtable(tx)
		if err != nil {
			return fmt.Errorf("err : %w", err)
		}

		err = db.InsertUserData(userID, userName, tx)
		if err != nil {
			return fmt.Errorf("err : %w", err)
		}
		return nil
	}, dbConn)
	if err != nil {
		return fmt.Errorf("err : %w", err)
	}
	return nil
}
