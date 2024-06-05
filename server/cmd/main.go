package main

import (
	"context"
	"errors"

	"github.com/0xdeafcafe/pillar-box/server/internals/broadcaster"
	"github.com/0xdeafcafe/pillar-box/server/internals/messagemonitor"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	log, err := zap.NewDevelopment(zap.AddStacktrace(zap.WarnLevel))
	if err != nil {
		panic(errors.Join(errors.New("failed to create logger"), errors.New(err.Error())))
	}

	broadcaster := broadcaster.New(ctx, log)
	broadcaster.ListenAndBroadcast()

	monitor, err := messagemonitor.New(ctx, log, broadcaster)
	if err != nil {
		panic(errors.Join(errors.New("failed to create monitor"), errors.New(err.Error())))
	}

	monitor.ListenAndHandle(ctx)
}
