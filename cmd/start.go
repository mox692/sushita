package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"time"

	"github.com/mox692/sushita/constant"
	"github.com/mox692/sushita/db"
	"golang.org/x/xerrors"

	"github.com/spf13/cobra"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "`start` starts sushita!",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		status := status()
		if status != 0 {
			return fmt.Errorf("status err: %d", status)
		}
		return nil
	},
}

func status() (exit exitcode.ExitCode) {
	exit = exitcode.Normal
	if err := start(); err != nil {
		fmt.Println(err)
		exit = exitcode.Abnormal
	}

	return exit
}

func start() error {

	user, err := user.Current()
	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	dbPath := user.HomeDir + "/db.sql"
	if f, err := os.Stat(dbPath); os.IsNotExist(err) || f.IsDir() {
		fmt.Printf("`db.sql` is not found in %s.\n Run `sushita init`.", dbPath)
		return nil
	}

	score := 0
	timer := time.NewTimer(time.Second * 15)

	s := bufio.NewScanner(os.Stdin)
	now_question := constant.DefaultWords[len(constant.DefaultWords)-1]
	fmt.Println(now_question)

	remainTime := 15
	go func() {
		ad := &remainTime
		for range time.Tick(1 * time.Second) {
			*ad = *ad - 1
		}

	}()

	go func() {
		for s.Scan() {
			if s.Text() == now_question {
				score++
				fmt.Printf("**********************************\n")
				fmt.Printf("collect!!\nTime Remain : %d\nScore : %d\n", remainTime, score)
				fmt.Printf("**********************************\n")
				now_question = constant.DefaultWords[rand.Intn(len(constant.DefaultWords))-1]
				fmt.Println(now_question)
			} else {
				fmt.Printf("incollect...\nTime Remain : %d\n\n", remainTime)
				fmt.Printf("%s\n", now_question)
			}
		}
		<-timer.C
	}()

	time.Sleep(time.Second * 15)
	fmt.Printf("\n\n")
	fmt.Printf("======================\n")
	fmt.Printf("time over!!\n your score : %d\n", score)
	fmt.Printf("======================\n")

	// もしハイスコアだったら、localdbにスコアを保存
	// dbからハイスコアを取得
	highScoreData, err := getHighScore()
	if err != nil {
		log.Fatal("err: %w", err)
	}
	if score > highScoreData.Score {
		err = askToSend(score)
	}
	if err != nil {
		log.Fatal("err: %w", err)
	}

	if s.Err() != nil {
		// non-EOF error.
		log.Fatal(s.Err())
		return s.Err()
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
	// ranking送信の確認テキスト
	fmt.Println("Do you want to send your highscore to the server? (Y/N)")
	// yes,noの受け取り
	var input string
	_, err := fmt.Scanf("%s", input)
	if err != nil {
		log.Fatal("err: %w", err)
		return err
	}

	switch input {
	case "Y", "y":
		err := sendRankingData(score)
		return err
	default:
		fmt.Println("Not sending.")
		return nil
	}
	return nil
}

func sendRankingData(score int) error {
	user, err := db.SelectUser()
	client := new(http.Client)
	url := "https://sushita.uc.r.appspot.com/ranking/set"
	req, err := http.NewRequest("Get", url, nil)
	req.Header.Set("user-token", strconv.Itoa(user.Id))
	res, err := client.Do(req)

	if err != nil {
		return fmt.Errorf(" %w", err)
	}

	fmt.Printf("%#v", res)
	return nil
}
