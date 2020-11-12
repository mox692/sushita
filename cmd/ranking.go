package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
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

type rankingDB interface {
	Select(key, value string) (*db.User, error)
}
type DB struct {
	DbConnection *sql.DB
}
type Ranking struct {
	endPoint   string
	httpMethod string
	user       db.User
}

func newDB() (*DB, error) {
	user, err := user.Current()
	if err != nil {
		return &DB{}, fmt.Errorf("err : %w", err)
	}
	dbPath := user.HomeDir + "/db.sql"
	DbConnection, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return &DB{}, fmt.Errorf("err : %w", err)
	}
	return &DB{DbConnection: DbConnection}, nil
}

func (d *DB) Select() (*db.User, error) {
	row := d.DbConnection.QueryRow("select * from user")
	return convertToUser(row)
}

type userScore struct {
	UserName string `json: "user_name"`
	Score    int    `json: "score"`
	Ranking  int    `json: "ranking"`
}

func getRanking() error {
	client := new(http.Client)
	req, err := http.NewRequest("GET", "https://sushita.uc.r.appspot.com/ranking", nil)
	user, _ := getUser()
	req.Header.Add("user-token", user.Id)
	resp, err := client.Do(req)
	if err != nil {
		return xerrors.Errorf("client.Do err : %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)

	var userScores []userScore
	err = json.Unmarshal(byteArray, &userScores)
	if err != nil {
		return xerrors.Errorf("json.Unmarshal err : %w", err)
	}
	t, _ := template.New("temp").Parse(constant.RankingLog)

	for _, v := range userScores {
		t.Execute(os.Stdout, v)
	}
	return nil
}

func getUser() (*db.User, error) {
	row := db.DbConnection.QueryRow("select * from user")
	return convertToUser(row)
}

func convertToUser(row *sql.Row) (*db.User, error) {
	u := db.User{}
	err := row.Scan(&u.Id, &u.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, xerrors.Errorf("row.Scan error: %w", err)
	}
	return &u, nil
}
