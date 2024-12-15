package app

import (
	"errors"

	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/broadcaster"
	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/macos"
	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/messagemonitor"
)

type App struct {
	Broadcaster *broadcaster.Broadcaster
	Monitor     *messagemonitor.MessageMonitor
	MacOS       *macos.MacOS
}

func New(debug bool) *App {
	monitor, err := messagemonitor.New()
	if err != nil {
		panic(errors.Join(errors.New("failed to create monitor"), err))
	}

	broadcaster := broadcaster.New()
	macos := macos.New(monitor, debug)

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
	a.Monitor.RegisterNoAccessHandler(a.MacOS.HandleNoAccess)

	// Run server and monitor in go routines
	go a.Broadcaster.ListenAndBroadcast()
	go a.Monitor.ListenAndHandle()

	a.MacOS.Run()
}
