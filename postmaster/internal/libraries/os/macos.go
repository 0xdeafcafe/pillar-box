package os

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/caseymrm/menuet"
	"golang.design/x/clipboard"

	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/messagemonitor"
	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/updater"
)

type UndefinedBool int

const (
	prefCopyCodeToClipboard  = "com.0xdeafcafe.pillar-box-postmaster_copy-code-to-clipboard"
	prefGetPrereleaseUpdates = "com.0xdeafcafe.pillar-box-postmaster_get-prerelease-updates"

	UndefinedBoolUndefined UndefinedBool = iota
	UndefinedBoolTrue
	UndefinedBoolFalse
)

type MacOS struct {
	debug   bool
	monitor *messagemonitor.MessageMonitor

	updater *updater.Updater

	latestCode  *MacOSLatestCode
	preferences *MacOSPreferences
}

type MacOSPreferences struct {
	CopyCodeToClipboard  bool
	GetPrereleaseUpdates bool
}

type MacOSLatestCode struct {
	DetectedAt time.Time
	MFACode    string
}

// New creates a new MacOS instance. The MacOS instance is responsible for managing the
// macOS menu bar application and rendering the menu items. The MacOS instance is also
// responsible for handling MFA codes detected by the MessageMonitor, displaying them
// in the menu, and copying them to the clipboard.
func NewMacOS(monitor *messagemonitor.MessageMonitor, debug bool) *MacOS {
	err := clipboard.Init()
	if err != nil {
		log.Printf("failed to initialize clipboard: %v", err)
	}

	macos := &MacOS{
		debug:   debug,
		monitor: monitor,

		updater: updater.New(),

		preferences: &MacOSPreferences{},
	}

	updater := updater.New()
	updater.RegisterNewVersionAvailableHandler(macos.HandleNewVersionAvailable)
	updater.RegisterGetPrereleasePreferenceHandler(func() bool {
		return macos.preferences.GetPrereleaseUpdates
	})

	return macos
}

func (m *MacOS) HandleMFACode(mfaCode string) {
	m.latestCode = &MacOSLatestCode{
		DetectedAt: time.Now(),
		MFACode:    mfaCode,
	}

	if m.preferences.CopyCodeToClipboard {
		clipboard.Write(clipboard.FmtText, []byte(mfaCode))

		menuet.App().Notification(menuet.Notification{
			Title:                        "New code detected",
			Subtitle:                     fmt.Sprintf("Code: %s", mfaCode),
			Message:                      "Copied to clipboard",
			RemoveFromNotificationCenter: true,
		})
	}

	m.renderMenu()
}

func (m *MacOS) HandleNoAccess() {
	response := menuet.App().Alert(menuet.Alert{
		MessageText:     "Pillar Box needs Full Disk Access to read incoming codes",
		InformativeText: "To enable Full Disk Access: System Settings > Security & Privacy > Full Disk Access > Add Pillar Box",
		Buttons:         []string{"OK", "Cancel"},
	})

	if response.Button == 0 {
		cmd := exec.Command("open", "x-apple.systempreferences:com.apple.preference.security?Privacy_AllFiles")

		log.Printf("opening system preferences: %v", cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("failed to open system preferences: %v", err)
		}
	}
}

func (m *MacOS) HandleNewVersionAvailable(name, version, url string) {
	response := menuet.App().Alert(menuet.Alert{
		MessageText:     "Update Available",
		InformativeText: "A new version of Pillar Box is available. Would you like to download it?",
		Buttons:         []string{"Download", "Later"},
	})
	if response.Button == 0 {
		cmd := exec.Command("open", url)

		log.Printf("opening browser: %v", cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("failed to open browser: %v", err)
		}
	}
}

func (m *MacOS) Run() {
	m.updater.StartBackgroundChecker()

	m.renderMenu()

	m.preferences.CopyCodeToClipboard = readAndSanitiseBoolPref(prefCopyCodeToClipboard)

	menuet.App().RunApplication()
}

func (m *MacOS) renderMenu() {
	menuet.App().Label = "com.0xdeafcafe.pillar-box-postmaster"
	menuet.App().Children = m.createMenuItems

	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "PB",
	})
}

func (m *MacOS) createMenuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{
		m.createLatestDiscoveredMenuItem(),
	}

	if m.debug {
		items = append(items, m.createDebugFakeMessageInitiatorMenuItem())
	}

	items = append(items,
		m.createCopyLastCodeMenuItem(),
		menuet.MenuItem{Type: menuet.Separator},
		m.createCopyCodesToClipboardMenuItem(),
		m.cretePrereleaseUpdatesMenuItem(),
	)

	return items
}

func (m *MacOS) createLatestDiscoveredMenuItem() menuet.MenuItem {
	if m.latestCode == nil {
		return menuet.MenuItem{
			Text: "Listening for codes...",
		}
	}

	return menuet.MenuItem{
		Text: fmt.Sprintf("Latest discovered code: %s", m.latestCode.MFACode),
	}
}

func (m *MacOS) createCopyLastCodeMenuItem() menuet.MenuItem {
	if m.latestCode == nil {
		return menuet.MenuItem{
			Text: "No code to copy",
		}
	}

	return menuet.MenuItem{
		Text: "Copy latest code to clipboard",
		Clicked: func() {
			err := clipboard.Init()
			if err != nil {
				fmt.Printf("failed to initialize clipboard: %v\n", err)
				return
			}

			clipboard.Write(clipboard.FmtText, []byte(m.latestCode.MFACode))
		},
	}
}

func (m *MacOS) createCopyCodesToClipboardMenuItem() menuet.MenuItem {
	return menuet.MenuItem{
		Text:  "Automatically copy to clipboard",
		State: m.preferences.CopyCodeToClipboard,
		Clicked: func() {
			newState := !m.preferences.CopyCodeToClipboard

			m.preferences.CopyCodeToClipboard = newState
			menuet.Defaults().SetBoolean(prefCopyCodeToClipboard, newState)
		},
	}
}

func (m *MacOS) cretePrereleaseUpdatesMenuItem() menuet.MenuItem {
	return menuet.MenuItem{
		Text:  "Get pre-release updates",
		State: m.preferences.GetPrereleaseUpdates,
		Clicked: func() {
			newState := !m.preferences.GetPrereleaseUpdates

			m.preferences.GetPrereleaseUpdates = newState
			menuet.Defaults().SetBoolean(prefGetPrereleaseUpdates, newState)

			if err := m.updater.CheckForUpdates(); err != nil {
				log.Printf("failed to check for updates: %v", err)
			}
		},
	}
}

func (m *MacOS) createDebugFakeMessageInitiatorMenuItem() menuet.MenuItem {
	return menuet.MenuItem{
		Text: "[debug] Dispatch random mock MFA code (5 second fuse)",
		Clicked: func() {
			m.monitor.SendMockMessage()
		},
	}
}

func readAndSanitiseBoolPref(pref string) bool {
	prefValue := UndefinedBool(menuet.Defaults().Integer(pref))

	if prefValue < UndefinedBoolUndefined || prefValue > UndefinedBoolFalse {
		return true
	}
	if prefValue == UndefinedBoolTrue {
		return true
	}

	return false
}
