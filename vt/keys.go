package vt

// Cursor Keys
const (
	KeyUp    = CSI + "A"
	KeyDown  = CSI + "B"
	KeyRight = CSI + "C"
	KeyLeft  = CSI + "D"
	KeyHome  = CSI + "H"
	KeyEnd   = CSI + "F"

	KeyCtrlUp    = CSI + "1;5A"
	KeyCtrlDown  = CSI + "1;5B"
	KeyCtrlRight = CSI + "1;5C"
	KeyCtrlLeft  = CSI + "1;5D"

	KeyAltUp    = CSI + "1;3A"
	KeyAltDown  = CSI + "1;3B"
	KeyAltRight = CSI + "1;3C"
	KeyAltLeft  = CSI + "1;3D"
)

// Numpad & Function Keys
const (
	KeyBackspace = "\x7f"
	KeyPause     = "\x1a"
	KeyEsc       = ESC
	KeyInsert    = CSI + "2~"
	KeyDelete    = CSI + "3~"
	KeyPageUp    = CSI + "5~"
	KeyPageDown  = CSI + "6~"

	KeyF1  = CSI + "11~"
	KeyF2  = CSI + "12~"
	KeyF3  = CSI + "13~"
	KeyF4  = CSI + "14~"
	KeyF5  = CSI + "15~"
	KeyF6  = CSI + "17~"
	KeyF7  = CSI + "18~"
	KeyF8  = CSI + "19~"
	KeyF9  = CSI + "20~"
	KeyF10 = CSI + "21~"
	KeyF11 = CSI + "23~"
	KeyF12 = CSI + "24~"
)
