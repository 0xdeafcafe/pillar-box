package macos

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/caseymrm/menuet"
	"golang.design/x/clipboard"

	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/messagemonitor"
)

type UndefinedBool int

const (
	prefCopyCodeToClipboard = "com.0xdeafcafe.pillar-box-postmaster_copy-code-to-clipboard"

	UndefinedBoolUndefined UndefinedBool = iota
	UndefinedBoolTrue
	UndefinedBoolFalse
)

type MacOS struct {
	debug   bool
	monitor *messagemonitor.MessageMonitor

	LatestCode  *MacOSLatestCode
	Preferences *MacOSPreferences
}

type MacOSPreferences struct {
	CopyCodeToClipboard bool
}

type MacOSLatestCode struct {
	DetectedAt time.Time
	MFACode    string
}

// New creates a new MacOS instance. The MacOS instance is responsible for managing the
// macOS menu bar application and rendering the menu items. The MacOS instance is also
// responsible for handling MFA codes detected by the MessageMonitor, displaying them
// in the menu, and copying them to the clipboard.
func New(monitor *messagemonitor.MessageMonitor, debug bool) *MacOS {
	err := clipboard.Init()
	if err != nil {
		log.Printf("failed to initialize clipboard: %v", err)
	}

	return &MacOS{
		debug:   debug,
		monitor: monitor,

		Preferences: &MacOSPreferences{},
	}
}

func (m *MacOS) HandleMFACode(mfaCode string) {
	m.LatestCode = &MacOSLatestCode{
		DetectedAt: time.Now(),
		MFACode:    mfaCode,
	}

	if m.Preferences.CopyCodeToClipboard {
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
	log.Printf("response.Button: %d", response.Button)

	if response.Button == 0 {
		cmd := exec.Command("open", "x-apple.systempreferences:com.apple.preference.security?Privacy_AllFiles")

		log.Printf("opening system preferences: %v", cmd)

		if err := cmd.Run(); err != nil {
			log.Printf("failed to open system preferences: %v", err)
		}
	}
}

func (m *MacOS) Run() {
	m.renderMenu()

	m.Preferences.CopyCodeToClipboard = readAndSanitiseBoolPref(prefCopyCodeToClipboard)

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
		menuet.MenuItem{Type: menuet.Separator},
		m.createCopyLastCodeMenuItem(),
		menuet.MenuItem{Type: menuet.Separator},
		m.createCopyCodesToClipboardMenuItem(),
	)

	return items
}

func (m *MacOS) createLatestDiscoveredMenuItem() menuet.MenuItem {
	if m.LatestCode == nil {
		return menuet.MenuItem{
			Text: "Listening for codes...",
		}
	}

	return menuet.MenuItem{
		Text: fmt.Sprintf("Latest discovered code: %s", m.LatestCode.MFACode),
	}
}

func (m *MacOS) createCopyLastCodeMenuItem() menuet.MenuItem {
	if m.LatestCode == nil {
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

			clipboard.Write(clipboard.FmtText, []byte(m.LatestCode.MFACode))
		},
	}
}

func (m *MacOS) createCopyCodesToClipboardMenuItem() menuet.MenuItem {
	return menuet.MenuItem{
		Text:  "Copy code to clipboard",
		State: m.Preferences.CopyCodeToClipboard,
		Clicked: func() {
			newState := !m.Preferences.CopyCodeToClipboard

			m.Preferences.CopyCodeToClipboard = newState
			menuet.Defaults().SetBoolean(prefCopyCodeToClipboard, newState)
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
