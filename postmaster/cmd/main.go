package main

import (
	"errors"
	"os"

	"go.uber.org/zap"

	"github.com/0xdeafcafe/pillar-box/server/internal/app"
)

func main() {
	var debug bool

	if len(os.Args) > 0 {
		debug = os.Args[1] == "--debug"
	}

	log, err := zap.NewDevelopment(zap.AddStacktrace(zap.WarnLevel))
	if err != nil {
		panic(errors.Join(errors.New("failed to create logger"), errors.New(err.Error())))
	}

	app := app.New(log, debug)
	app.Run()
}
