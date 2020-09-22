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
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// `sushita start` を実行した時のステータスコードのテスト
func TestStartSushita(t *testing.T) {
	status := status()
	fmt.Println(status, "status!!")
	if status != 0 {
		t.Errorf("expected 0, get %d", status)
	}
}

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
