package main

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"

	"github.com/mox692/sushita/cmd"
)

func main() {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())
	cmd.Execute()
}
