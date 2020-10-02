package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// cmdパッケージのExecuteの終了ステータスコードのテスト
func TestExecute(t *testing.T) {
	status := Execute()
	fmt.Println(status, "status!!")
	if status != 0 {
		t.Errorf("expected 0, get %d", status)
	}
}

// 期待した標準出力が得られるかのテスト
func TestPrint(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = w

	Execute()

	os.Stdout = stdout
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	if strings.TrimRight(buf.String(), "\r\n") != "Let's run `sushita start`" {
		t.Errorf("Execute() = %s, want 'Let's run `sushita start`'", buf.String())
	}
}
