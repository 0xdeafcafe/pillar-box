package app

import (
	"errors"

	"github.com/0xdeafcafe/pillar-box/server/internals/broadcaster"
	"github.com/0xdeafcafe/pillar-box/server/internals/macos"
	"github.com/0xdeafcafe/pillar-box/server/internals/messagemonitor"
	"go.uber.org/zap"
)

type App struct {
	Broadcaster *broadcaster.Broadcaster
	Monitor     *messagemonitor.MessageMonitor
	MacOS       *macos.MacOS
}

func New(log *zap.Logger) *App {
	broadcaster := broadcaster.New(log)
	macos := macos.New()

	monitor, err := messagemonitor.New(log)
	if err != nil {
		panic(errors.Join(errors.New("failed to create monitor"), errors.New(err.Error())))
	}

	return &App{
		Broadcaster: broadcaster,
		Monitor:     monitor,
		MacOS:       macos,
	}
}

func (a *App) Run() {
	// setup detection handlers
	a.Monitor.SetDetectionHandler(func(mfaCode string) {
		a.Broadcaster.BroadcastMFACode(mfaCode)
		a.MacOS.HandleMFACode(mfaCode)
	})

	go a.Broadcaster.ListenAndBroadcast()
	go a.Monitor.ListenAndHandle()
	a.MacOS.Run()
}
