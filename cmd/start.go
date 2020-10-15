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

	/*
		1).15秒の時間を管理するgoroutine
		2).問題をひたすら表示させるgoroutine

		・15sのチャネルを常に監視する(switchする)処理を2)に入れる。
		・time sleepはダサいので、sync.Waitを調べて導入する。
	*/

	// [goroutine1]
	// バックグラウンドでタイマー値の減数を処理する。
	// Todo: timeパッケージで、もっとマシな便利メソッド的なのはないか？？

	wg := sync.WaitGroup{}
	timeover := make(chan struct{})
	wg.Add(1)
	go func() {
		ad := &remainTime

		// 1s毎にremain Timeを1つずつデクリメント
		for range time.Tick(1 * time.Second) {
			*ad = *ad - 1
			// timer終了
			if *ad == 0 {
				timer.Stop()
				// timer.Cをしっかり読み捨てる
				timeover <- struct{}{}
				fmt.Println("timer終了")
				break
			}
		}
		wg.Done()
		fmt.Println("wgdone 終わり")
	}()

	// [goroutine2]
	// 無限に常に標準入力を待つ

	wg.Add(1)
	go func() {
		// 改行毎にscan
		for s.Scan() {
			select {
			case <-timeover:
				goto L
			// ここをdefaultジャなくて、入力が来たことを通知するチャネルを作って
			// それに変更できないか

			default:
				if s.Text() == now_question {
					score++
					// Todo: ここを1つで表現
					fmt.Printf("**********************************\n")
					fmt.Printf("collect!!\nTime Remain : %d\nScore : %d\n", remainTime, score)
					fmt.Printf("**********************************\n")
					// ここまで
					now_question = constant.DefaultWords[rand.Intn(len(constant.DefaultWords))-1]
					fmt.Println(now_question)
				} else {
					fmt.Printf("incollect...\nTime Remain : %d\n\n", remainTime)
					fmt.Printf("%s\n", now_question)
				}
			}
			// Todo: ここをtimer.Cのswitchにできないか。
		}
	L:
		wg.Done()
		fmt.Println("task2 done")
	}()

	fmt.Println("wating....")
	wg.Wait()
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
	fmt.Println("high score:", highScoreData.Score)
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

	fmt.Printf("\n\n🎉🎉🎉= HIGH SCORE !!! =🎉🎉🎉\n\n")
	fmt.Println("Do you want to send your highscore to the server? (Y/N)")

	var input string
	_, err := fmt.Scanf("%s", &input)
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
		return fmt.Errorf(" %w", err)
	}

	// Todo: レスポンスの構造体を決めて、decode処理を実装
	// fmt.Printf("%#v\n\n", res)
	// bodyを表示
	//	***jsonにデコードする時***
	// json.NewDecoder(res.Body).Decode(res)
	// buf, _ := json.Marshal(res.Body)
	b, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))
	return nil
}

type sendRankingRequest struct {
	Name  string `json: "name"`
	Score int    `json: "score"`
}

type response struct {
	Score int32 `json:"score"`
}
