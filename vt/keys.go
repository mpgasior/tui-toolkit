package vt

//go:generate stringer -type=Key
type Key int

const (
	KeyUnknown Key = iota

	KeyTab
	KeyShiftTab

	KeyUp
	KeyDown
	KeyRight
	KeyLeft

	KeyHome
	KeyEnd

	KeyCtrlUp
	KeyCtrlDown
	KeyCtrlRight
	KeyCtrlLeft

	KeyAltUp
	KeyAltDown
	KeyAltRight
	KeyAltLeft

	KeyShiftUp
	KeyShiftDown
	KeyShiftRight
	KeyShiftLeft

	KeyBackspace
	KeyPause
	KeyEsc
	KeyInsert
	KeyDelete
	KeyPageUp
	KeyPageDown

	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12

	KeyCtrlA
	KeyCtrlB
	KeyCtrlC
	KeyCtrlD
	KeyCtrlE
	KeyCtrlF
	KeyCtrlG
	KeyCtrlH
	KeyCtrlI
	KeyCtrlJ
	KeyCtrlK
	KeyCtrlL
	KeyCtrlM
	KeyCtrlN
	KeyCtrlO
	KeyCtrlP
	KeyCtrlQ
	KeyCtrlR
	KeyCtrlS
	KeyCtrlT
	KeyCtrlU
	KeyCtrlV
	KeyCtrlW
	KeyCtrlX
	KeyCtrlY
	KeyCtrlZ

	KeyAltA
	KeyAltB
	KeyAltC
	KeyAltD
	KeyAltE
	KeyAltF
	KeyAltG
	KeyAltH
	KeyAltI
	KeyAltJ
	KeyAltK
	KeyAltL
	KeyAltM
	KeyAltN
	KeyAltO
	KeyAltP
	KeyAltQ
	KeyAltR
	KeyAltS
	KeyAltT
	KeyAltU
	KeyAltV
	KeyAltW
	KeyAltX
	KeyAltY
	KeyAltZ
)

var KeySequences = map[Key]string{
	KeyUnknown:  "",
	KeyTab:      "\t",
	KeyShiftTab: CSI + "Z",

	KeyUp:    CSI + "A",
	KeyDown:  CSI + "B",
	KeyRight: CSI + "C",
	KeyLeft:  CSI + "D",

	KeyHome: CSI + "H",
	KeyEnd:  CSI + "F",

	KeyCtrlUp:    CSI + "1;5A",
	KeyCtrlDown:  CSI + "1;5B",
	KeyCtrlRight: CSI + "1;5C",
	KeyCtrlLeft:  CSI + "1;5D",

	KeyAltUp:    CSI + "1;3A",
	KeyAltDown:  CSI + "1;3B",
	KeyAltRight: CSI + "1;3C",
	KeyAltLeft:  CSI + "1;3D",

	KeyShiftUp:    CSI + "1;2A",
	KeyShiftDown:  CSI + "1;2B",
	KeyShiftRight: CSI + "1;2C",
	KeyShiftLeft:  CSI + "1;2D",

	KeyBackspace: "\x7f",
	KeyPause:     "\x1a",
	KeyEsc:       ESC,
	KeyInsert:    CSI + "2~",
	KeyDelete:    CSI + "3~",
	KeyPageUp:    CSI + "5~",
	KeyPageDown:  CSI + "6~",

	KeyF1:  ESC + "OP",
	KeyF2:  ESC + "OQ",
	KeyF3:  ESC + "OR",
	KeyF4:  ESC + "OS",
	KeyF5:  CSI + "15~",
	KeyF6:  CSI + "17~",
	KeyF7:  CSI + "18~",
	KeyF8:  CSI + "19~",
	KeyF9:  CSI + "20~",
	KeyF10: CSI + "21~",
	KeyF11: CSI + "23~",
	KeyF12: CSI + "24~",

	KeyCtrlA: "\x01",
	KeyCtrlB: "\x02",
	KeyCtrlC: "\x03",
	KeyCtrlD: "\x04",
	KeyCtrlE: "\x05",
	KeyCtrlF: "\x06",
	KeyCtrlG: "\x07",
	KeyCtrlH: "\x08",
	KeyCtrlI: "\x09",
	KeyCtrlJ: "\x0a",
	KeyCtrlK: "\x0b",
	KeyCtrlL: "\x0c",
	KeyCtrlM: "\x0d",
	KeyCtrlN: "\x0e",
	KeyCtrlO: "\x0f",
	KeyCtrlP: "\x10",
	KeyCtrlQ: "\x11",
	KeyCtrlR: "\x12",
	KeyCtrlS: "\x13",
	KeyCtrlT: "\x14",
	KeyCtrlU: "\x15",
	KeyCtrlV: "\x16",
	KeyCtrlW: "\x17",
	KeyCtrlX: "\x18",
	KeyCtrlY: "\x19",
	KeyCtrlZ: "\x1a",

	KeyAltA: ESC + "a",
	KeyAltB: ESC + "b",
	KeyAltC: ESC + "c",
	KeyAltD: ESC + "d",
	KeyAltE: ESC + "e",
	KeyAltF: ESC + "f",
	KeyAltG: ESC + "g",
	KeyAltH: ESC + "h",
	KeyAltI: ESC + "i",
	KeyAltJ: ESC + "j",
	KeyAltK: ESC + "k",
	KeyAltL: ESC + "l",
	KeyAltM: ESC + "m",
	KeyAltN: ESC + "n",
	KeyAltO: ESC + "o",
	KeyAltP: ESC + "p",
	KeyAltQ: ESC + "q",
	KeyAltR: ESC + "r",
	KeyAltS: ESC + "s",
	KeyAltT: ESC + "t",
	KeyAltU: ESC + "u",
	KeyAltV: ESC + "v",
	KeyAltW: ESC + "w",
	KeyAltX: ESC + "x",
	KeyAltY: ESC + "y",
	KeyAltZ: ESC + "z",
}
