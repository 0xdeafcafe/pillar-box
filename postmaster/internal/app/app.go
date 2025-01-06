package app

import (
	"errors"

	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/broadcaster"
	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/messagemonitor"
	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/os"
	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/updater"
)

type App struct {
	Broadcaster *broadcaster.Broadcaster
	Monitor     *messagemonitor.MessageMonitor
	OS          os.OS
	Updater     *updater.Updater
}

func New(debug bool) *App {
	monitor, err := messagemonitor.New()
	if err != nil {
		panic(errors.Join(errors.New("failed to create monitor"), err))
	}

	broadcaster := broadcaster.New()

	os, err := os.New(monitor, debug)
	if err != nil {
		panic(errors.Join(errors.New("failed to create OS"), err))
	}

	return &App{
		Broadcaster: broadcaster,
		Monitor:     monitor,
		OS:          os,
		Updater:     updater.New(),
	}
}

func (a *App) Run() {
	// Setup detection handlers
	a.Monitor.RegisterDetectionHandler(a.Broadcaster.BroadcastMFACode)
	a.Monitor.RegisterDetectionHandler(a.OS.HandleMFACode)
	a.Monitor.RegisterNoAccessHandler(a.OS.HandleNoAccess)

	// Setup update handlers
	a.Updater.RegisterNewVersionAvailableHandler(a.OS.HandleNewVersionAvailable)

	// Run server and monitor in go routines
	go a.Broadcaster.ListenAndBroadcast()
	go a.Monitor.ListenAndHandle()
	a.Updater.StartBackgroundChecker()

	a.OS.Run()
}
