package macos

import (
	"fmt"
	"time"

	"github.com/caseymrm/menuet"
	"golang.design/x/clipboard"

	"github.com/0xdeafcafe/pillar-box/server/internal/libraries/messagemonitor"
)

const (
	prefCopyCodes = "com.0xdeafcafe.pillar-box-postmaster_also-copy-codes-to-clipboard"
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

	m.renderMenu()
}

func (m *MacOS) Run() {
	m.renderMenu()

	m.Preferences.CopyCodeToClipboard = menuet.Defaults().Boolean(prefCopyCodes)

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
		m.createCopyLastCodeMenuItem(),
		m.createCopyCodesToClipboardMenuItem(),
	}

	if m.debug {
		items = append(items, m.createDebugFakeMessageInitiatorMenuItem())
	}

	return items
}

func (m *MacOS) createLatestDiscoveredMenuItem() menuet.MenuItem {
	if m.LatestCode == nil {
		return menuet.MenuItem{
			Text: "Waiting for codes...",
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
		Text:  "Also copy codes to clipboard",
		State: m.Preferences.CopyCodeToClipboard,
		Clicked: func() {
			newState := !m.Preferences.CopyCodeToClipboard

			m.Preferences.CopyCodeToClipboard = newState
			menuet.Defaults().SetBoolean(prefCopyCodes, newState)
		},
	}
}

func (m *MacOS) createDebugFakeMessageInitiatorMenuItem() menuet.MenuItem {
	return menuet.MenuItem{
		Text: "[debug] Dispatch random mock 2FA code (5 second fuse)",
		Clicked: func() {
			newState := !m.Preferences.CopyCodeToClipboard

			m.Preferences.CopyCodeToClipboard = newState
			menuet.Defaults().SetBoolean(prefCopyCodes, newState)
		},
	}
}
