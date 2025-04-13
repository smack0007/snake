package main

import (
	"os"

	"github.com/smack0007/snake/engine"
)

func main() {
	os.Exit(run())
}

func run() int {
	game := NewGame()
	err := engine.Run(&game)

	if err != nil {
		return 1
	}

	return 0
}
