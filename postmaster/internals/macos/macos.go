package macos

import (
	"fmt"
	"time"

	"github.com/caseymrm/menuet"
	"golang.design/x/clipboard"
)

const (
	prefCopyCodes = "com.0xdeafcafe.pillar-box-postmaster_also-copy-codes-to-clipboard"
)

type MacOS struct {
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

func New() *MacOS {
	return &MacOS{
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

	// load preferences
	m.Preferences.CopyCodeToClipboard = menuet.Defaults().Boolean(prefCopyCodes)

	menuet.App().RunApplication()
}

func (m *MacOS) renderMenu() {
	menuet.App().Label = "com.0xdeafcafe.pillar-box-postmaster"
	menuet.App().Children = m.generateMenuItems

	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "PB",
	})
}

func (m *MacOS) generateMenuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{
		m.generateLatestDiscoveredMenuItem(),
		m.generateCopyLastCodeMenuItem(),
		m.generateCopyCodesToClipboardMenuItem(),
	}

	return items
}

func (m *MacOS) generateLatestDiscoveredMenuItem() menuet.MenuItem {
	if m.LatestCode == nil {
		return menuet.MenuItem{
			Text: "Waiting for codes...",
		}
	}

	return menuet.MenuItem{
		Text: fmt.Sprintf("Latest discovered code: %s", m.LatestCode.MFACode),
	}
}

func (m *MacOS) generateCopyLastCodeMenuItem() menuet.MenuItem {
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

func (m *MacOS) generateCopyCodesToClipboardMenuItem() menuet.MenuItem {
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
