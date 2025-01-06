package main

import (
	"os"

	"github.com/0xdeafcafe/pillar-box/server/internal/app"
)

func main() {
	var debug bool
	if len(os.Args) > 1 {
		debug = os.Args[1] == "--debug"
	}

	app := app.New(debug)
	app.Run()
}
