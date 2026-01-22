package vt

// Query State
// See: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#query-state
const (
	QueryCursorPosition   = CSI + "6n"
	QueryDeviceAttributes = CSI + "0c"
	QueryTerminalName     = CSI + ">q"

	QueryOSCFmt      = OSC + "%d;?" + BEL
	QueryFgColor     = OSC + "10;?" + BEL
	QueryBgColor     = OSC + "11;?" + BEL
	QueryCursorColor = OSC + "12;?" + BEL
)
