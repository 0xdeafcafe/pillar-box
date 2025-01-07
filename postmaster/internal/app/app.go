package app

import (
	"errors"

	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/broadcaster"
	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/messagemonitor"
	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/os"
)

type App struct {
	Broadcaster *broadcaster.Broadcaster
	Monitor     *messagemonitor.MessageMonitor
	OS          os.OS
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
	}
}

func (a *App) Run() {
	// Setup detection handlers
	a.Monitor.RegisterDetectionHandler(a.Broadcaster.BroadcastMFACode)
	a.Monitor.RegisterDetectionHandler(a.OS.HandleMFACode)
	a.Monitor.RegisterNoAccessHandler(a.OS.HandleNoAccess)

	// Run server and monitor in go routines
	go a.Broadcaster.ListenAndBroadcast()
	go a.Monitor.ListenAndHandle()

	a.OS.Run()
}
