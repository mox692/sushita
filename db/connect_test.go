package db

import (
	"testing"
)

func TestDBinit(t *testing.T) {
	// DBinit()でerrが発生しないか
	if err := DBinit(); err != nil {
		t.Fatalf("err : %s", err)
	}
}
