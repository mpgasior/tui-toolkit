package vt

// Text Modification
// See: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#text-modification
const (
	InsertCharFmt = CSI + "%d@"
	DeleteCharFmt = CSI + "%dP"
	EraseCharFmt  = CSI + "%dX"
	InsertLineFmt = CSI + "%dL"
	DeleteLineFmt = CSI + "%dM"

	EraseLineFromCursor = CSI + "0K"
	EraseLineFromStart  = CSI + "1K"
	EraseEntireLine     = CSI + "2K"

	EraseScreenToBottom = CSI + "0J"
	EraseEntireScreen   = CSI + "2J"
)
