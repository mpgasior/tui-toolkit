package vt

// Query State
// See: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#query-state
const (
	ReportCursorPosition   = CSI + "6n"
	ReportDeviceAttributes = CSI + "0c"
)
