package app

import (
	"errors"

	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/broadcaster"
	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/macos"
	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/messagemonitor"
	"go.uber.org/zap"
)

type App struct {
	Broadcaster *broadcaster.Broadcaster
	Monitor     *messagemonitor.MessageMonitor
	MacOS       *macos.MacOS
}

func New(log *zap.Logger) *App {
	monitor, err := messagemonitor.New(log)
	if err != nil {
		panic(errors.Join(errors.New("failed to create monitor"), err))
	}

	broadcaster := broadcaster.New(log)
	macos := macos.New()

	return &App{
		Broadcaster: broadcaster,
		Monitor:     monitor,
		MacOS:       macos,
	}
}

func (a *App) Run() {
	// setup detection handlers
	a.Monitor.RegisterDetectionHandler(a.Broadcaster.BroadcastMFACode)
	a.Monitor.RegisterDetectionHandler(a.MacOS.HandleMFACode)

	// Run server and monitor in go routines
	go a.Broadcaster.ListenAndBroadcast()
	go a.Monitor.ListenAndHandle()

	a.MacOS.Run()
}
