package vt

import (
	"fmt"
	"io"
	"sync"
)

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

func EnterMode(w io.Writer, m Mode) (restore func() error, err error) {
	_, err = io.WriteString(w, m.Enable())

	var once sync.Once
	var restoreErr error
	restore = func() error {
		once.Do(func() {
			_, restoreErr = io.WriteString(w, m.Disable())
		})

		return restoreErr
	}

	return restore, err
}

// Mode Changes
// See: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#mode-changes
const (
	KeypadApplicationMode = ESC + "="
	KeypadNumericMode     = ESC + ">"
)
