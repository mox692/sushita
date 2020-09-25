package db

import (
	"os"
	"os/user"
	"testing"
)

func TestDBinit(t *testing.T) {
	// DBinit()でerrが発生しないか。
	if err := DBinit(); err != nil {
		t.Fatalf("err : %s", err)
	}

	// Todo: homedirにdb.sqファイルが存在するか。
	// テストだと、ファイルが作成されない。。
	user, err := user.Current()
	if err != nil {
		t.Fatal("user.Current() err ")
	}
	dbPath := user.HomeDir + "/db.sql"

	if f, err := os.Stat(dbPath); os.IsNotExist(err) || f.IsDir() {
		t.Fatalf("%s", err)
	}

}
