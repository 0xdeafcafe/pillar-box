package os

import (
	"fmt"
	"runtime"

	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/messagemonitor"
)

type OS interface {
	HandleMFACode(mfaCode string)
	HandleNoAccess()
	HandleNewVersionAvailable(name, version, url string)
	Run()
}

func New(monitor *messagemonitor.MessageMonitor, debug bool) (OS, error) {
	if runtime.GOOS == "darwin" {
		return NewMacOS(monitor, debug), nil
	}

	return nil, fmt.Errorf("unsupported OS: %s, only darwin is supported", runtime.GOOS)
}
