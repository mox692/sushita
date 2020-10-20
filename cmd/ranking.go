package cmd

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/mox692/sushita/constant"
	"github.com/mox692/sushita/db"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
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

func getRanking() error {
	client := new(http.Client)
	req, err := http.NewRequest("GET", "https://sushita.uc.r.appspot.com/ranking", nil)
	usr, _ := getUser()
	req.Header.Add("user-token", usr.UserId)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)

	var userScores []userScore
	err = json.Unmarshal(byteArray, &userScores)
	if err != nil {
		log.Fatal(err)
		return err
	}

	t, _ := template.New("temp").Parse(constant.RankingLog)

	for _, v := range userScores {
		t.Execute(os.Stdout, v)
	}

	return nil
}

type usr struct {
	UserId string
	Name   string
}

type userScore struct {
	Name    string `json:"name"`
	Score   int    `json:"score"`
	Ranking int    `json:ranking"`
}

func getUser() (*usr, error) {
	row := db.DbConnection.QueryRow("select * from user")
	return convertToUser(row)
}

func convertToUser(row *sql.Row) (*usr, error) {
	u := usr{}
	err := row.Scan(&u.UserId, &u.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, xerrors.Errorf("row.Scan error: %w", err)
	}
	return &u, nil
}
