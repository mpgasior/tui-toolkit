package vt

// Mode Changes
// See: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#mode-changes
const (
	UseKeypadApplicationMode = ESC + "="
	UseKeypadNumericMode     = ESC + ">"

	UseCursorKeysApplicationMode = CSI + "?1h"
	UseCursorKeysNumericMode     = CSI + "?1l"
)
