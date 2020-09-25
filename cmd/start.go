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
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/mox692/sushita/constant"

	"github.com/spf13/cobra"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		status()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

}

// この関数を挟むことでstatusコードのテストを行いやすくする。
func status() (exit exitcode.ExitCode) {
	exit = exitcode.Normal
	if err := startSushita(); err != nil {
		fmt.Println(err)
		exit = exitcode.Abnormal
	}
	return exit
}

// sushita startのメイン処理。
func startSushita() error {
	fmt.Println("↓enter eny command")
	score := 0
	timer := time.NewTimer(time.Second * 15)
	// question := []string{"きょうのばんごはん", "かのじょのゆくえ", "あしたのてんき"}

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

	if s.Err() != nil {
		// non-EOF error.
		log.Fatal(s.Err())
		return s.Err()
	}
	return nil
}
