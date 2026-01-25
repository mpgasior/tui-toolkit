package vt

import "fmt"

type Mode int

func (m Mode) Enable() string {
	return fmt.Sprintf("%s?%dh", CSI, m)
}

func (m Mode) Disable() string {
	return fmt.Sprintf("%s?%dl", CSI, m)
}

const (
	ModeCursorKeys      Mode = 1
	ModeShowCursor      Mode = 25
	ModeMouseAll        Mode = 1003
	ModeAlternateScreen Mode = 1049
	ModeBracketedPaste  Mode = 2004
)

// Mode Changes
// See: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#mode-changes
const (
	KeypadApplicationMode = ESC + "="
	KeypadNumericMode     = ESC + ">"
)
