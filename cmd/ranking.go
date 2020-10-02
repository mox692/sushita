package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// rankingCmd represents the ranking command
var rankingCmd = &cobra.Command{
	Use:   "ranking",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return getRanking()
	},
}

func init() {
	rootCmd.AddCommand(rankingCmd)

}

func getRanking() error {

	// serverにリクエストを送る処理
	client := new(http.Client)
	req, err := http.NewRequest("GET", "http://localhost:8080/ranking", nil)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray)) // htmlをstringで取得

	return nil
}
