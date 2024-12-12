package main

import (
	"errors"

	"go.uber.org/zap"

	"github.com/0xdeafcafe/pillar-box/server/internal/app"
)

func main() {
	log, err := zap.NewDevelopment(zap.AddStacktrace(zap.WarnLevel))
	if err != nil {
		panic(errors.Join(errors.New("failed to create logger"), errors.New(err.Error())))
	}

	app := app.New(log)
	app.Run()
}
