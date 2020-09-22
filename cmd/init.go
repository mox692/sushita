/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"../db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains example`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var userName string
		fmt.Printf("Enter your name!\n")
		fmt.Printf("=>")

		fmt.Scanf("%s", &userName)

		fmt.Printf("hello %s!!\n", userName)

		err := db.DBinit()
		if err != nil {
			fmt.Errorf("DbConnection.Exec : %w", err)
		}

		err = db.CreateUsertable(db.DbConnection)
		if err != nil {
			return fmt.Errorf("err : %w", err)
		}

		err = db.CreateRankingtable(db.DbConnection)
		if err != nil {
			return fmt.Errorf("err : %w", err)
		}

		cmd.Printf("init sushita!🎉\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
