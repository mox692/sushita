package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// `sushita start` を実行した時の標準出力のテスト
func TestPrintStart(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = w
	// !標準出力には表示されない
	status()

	os.Stdout = stdout
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	if !strings.Contains(strings.TrimRight(buf.String(), "\r\n"), "enter eny command") {
		t.Errorf("status() = %s, want 'enter eny command'", buf.String())
	}
}

// 正解した時の処理のテスト
func TestCollectAnswear(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("ここまで1")
	stdout := os.Stdout
	os.Stdout = w

	// !標準出力には表示されない
	status()
	w.Write([]byte("test"))
	// Todo: status中にwriteしたい。根本的に別の手法が必要？

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// 標準出力(os.Stdout)に元に戻す
	os.Stdout = stdout

	// 問題文を取得
	fmt.Println("ここまで2")
	// lastLine, err := getNthLine(&buf, 2)
	// fmt.Printf("%s", lastLine)

	if !strings.Contains(strings.TrimRight(buf.String(), "\r\n"), "enter eny command") {
		t.Errorf("status() = %s, want 'enter eny command'", buf.String())
	}
}

// 前処理と後処理で、
// 1. db.sqlがある時
// 2. ない時のテストをそれぞれ行う。
func TestStart_hasSQLFile(t *testing.T) {

}

// gameインスタンスをtableで用意して、tabledrivenテストを行う。
func TestStart_runGame(t *testing.T) {
}

func TestStart_insertGameScore(t *testing.T) {

}

// 標準入力のパターンをtableで用意してテスト(標準n入力をどのようにテストするか)
// Todo: sendRankingDataまで呼ばず、ここまで呼ばれたら成功！という仕組みにしたい。。
func TestStart_askToSend(t *testing.T) {

}

func TestStart_sendRankingData(t *testing.T) {

}

func getNthLine(text io.Reader, n int) (string, error) {
	line := ""
	scanner := bufio.NewScanner(text)
	// i := 1
	for scanner.Scan() {
		line = scanner.Text()
		fmt.Println("1行scan終了:", line)
		// if i == n {
		// 	break
		// }
		// i++
	}
	return line, nil
}
