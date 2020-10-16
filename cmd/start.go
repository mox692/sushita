package cmd

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/user"
	"sync"
	"time"

	"github.com/mox692/sushita/constant"
	"github.com/mox692/sushita/db"
	"golang.org/x/xerrors"

	"github.com/spf13/cobra"
)

type Game struct {
	question    []string
	nowQuestion string
	gameTime    time.Duration
	score       int
}

func newGame() *Game {
	return &Game{
		question:    constant.DefaultWords,
		nowQuestion: constant.DefaultWords[rand.Intn(len(constant.DefaultWords))],
		gameTime:    constant.InGameTime,
		score:       0,
	}
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "`start` starts sushita!",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return start()
	},
}

func start() error {

	err := hasSQLFile()
	if err != nil {
		return xerrors.Errorf("hasSqlFile err : %w", err)
	}

	game := newGame()
	fmt.Println(game.nowQuestion)
	s := bufio.NewScanner(os.Stdin)

	wg := sync.WaitGroup{}
	timeover := make(chan struct{})
	inputAnswer := make(chan string)

	wg.Add(1)
	go func() {
		for range time.Tick(1 * time.Second) {
			game.gameTime--
			if game.gameTime == 0 {
				timeover <- struct{}{}
				break
			}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for {
			select {
			case <-timeover:
				goto L
			case text := <-inputAnswer:
				if text == game.nowQuestion {
					game.score++
					fmt.Printf("**********************************\n")
					fmt.Printf("collect!!\nTime Remain : %d\nScore : %d\n", game.gameTime, game.score)
					fmt.Printf("**********************************\n")
					game.nowQuestion = constant.DefaultWords[rand.Intn(len(constant.DefaultWords))]
					fmt.Println(game.nowQuestion)
				} else {
					fmt.Printf("incollect...\nTime Remain : %d\n\n", game.gameTime)
					fmt.Printf("%s\n", game.nowQuestion)
				}
			}
		}
	L:
		wg.Done()
	}()

	go func() {
		for s.Scan() {
			inputAnswer <- s.Text()
		}
	}()

	wg.Wait()

	highScoreData, err := getHighScore()
	if err != nil {
		return xerrors.Errorf("getHighScore err : %w", err)
	}
	if highScoreData == nil {
		highScoreData = &db.LocalRanking{
			Score: 0,
		}
	}

	fmt.Printf("\n\n")
	fmt.Printf("======================\n")
	fmt.Printf("time over!!\n your score : %d\n high score: %d\n", game.score, highScoreData.Score)
	fmt.Printf("======================\n")

	// localにscoreを保存
	err = insertGameScore(game.score)
	if err != nil {
		return xerrors.Errorf("insertGameScore err : %w", err)
	}

	if game.score > highScoreData.Score {
		err = askToSend(game.score)
	}

	if err != nil {
		return xerrors.Errorf("askToSend err : %w", err)
	}

	if s.Err() != nil {
		// non-EOF error.
		log.Fatal(s.Err())
		return s.Err()
	}
	return nil
}

func hasSQLFile() error {
	user, err := user.Current()
	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	dbPath := user.HomeDir + "/db.sql"
	if f, err := os.Stat(dbPath); os.IsNotExist(err) || f.IsDir() {
		fmt.Printf("`db.sql` is not found in %s.\n Run `sushita init`.", dbPath)
		return nil
	}
	return nil
}

func insertGameScore(score int) error {
	stmt, err := db.DbConnection.Prepare("INSERT INTO local_ranking (score, created_at) VALUES (?, ?);")
	if err != nil {
		return xerrors.Errorf("db.Conn.Prepare err : %w", err)
	}
	_, err = stmt.Exec(score, timeToString(time.Now()))
	if err != nil {
		return xerrors.Errorf("stmt.Exec err : %w", err)
	}
	return nil
}

func getHighScore() (*db.LocalRanking, error) {
	localRanking := db.LocalRanking{}
	row := db.DbConnection.QueryRow("select * from local_ranking order by score desc limit 1;")
	err := row.Scan(&localRanking.Score, &localRanking.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, xerrors.Errorf("row.Scan error: %w", err)
	}
	return &localRanking, nil
}

func askToSend(score int) error {

	fmt.Printf("\n\n🎉🎉🎉= HIGH SCORE !!! =🎉🎉🎉\n\n")
	fmt.Println("Do you want to send your highscore to the server? (Y/N)")

	var inputAnswer string
	_, err := fmt.Scanf("%s", &inputAnswer)
	if err != nil {
		log.Fatal("err: %w", err)
		return err
	}

	switch inputAnswer {
	case "Y", "y":
		err := sendRankingData(score)
		return err
	default:
		fmt.Println("Not sending.")
		return nil
	}
}

func sendRankingData(score int) error {
	user, err := db.SelectUser()
	if err != nil {
		fmt.Errorf("err: %w", err)
	}
	client := new(http.Client)
	url := "https://sushita.uc.r.appspot.com/ranking/set"

	sendData := &sendRankingRequest{
		Name:  user.UserName,
		Score: score,
	}

	jsonData, err := json.Marshal(sendData)
	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	req, err := http.NewRequest("Get", url, bytes.NewBuffer(jsonData))
	req.Header.Set("user-token", user.Id)
	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		return xerrors.Errorf("client.Do err : %w", err)
	}

	b, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))
	fmt.Println("raknking dataを送信しました")

	return nil
}

func timeToString(t time.Time) string {
	var layout = "2006-01-02 15:04:05"
	return t.Format(layout)
}

type sendRankingRequest struct {
	Name  string `json: "name"`
	Score int    `json: "score"`
}

type response struct {
	Score int32 `json:"score"`
}
