package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
)

var rootCmd = &cobra.Command{
	Use:   "sushita",
	Short: "sushita command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Let's run `sushita start`")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(rankingCmd)
}

func Execute() (exit exitcode.ExitCode) {
	exit = exitcode.Normal
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		exit = exitcode.Abnormal
	}
	return exit
}
