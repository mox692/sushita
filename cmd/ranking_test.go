package cmd

import "testing"

// 現状のgetuserの実装がdbに依存したものになっている。
// getuser関数をモック化する。

func TestgetRanking(t *testing.T) {
	db, err := newDB()
	if err != nil {
		t.Fatalf("newDB err")
	}

}

//getuserのテスト
func TestgetUser(t *testing.T) {
	db, err := newDB()

}
