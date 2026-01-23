package vt

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

	KeyEnter
	KeyBackspace
	KeyPause
	KeyEsc
	KeyInsert
	KeyDelete
	KeyPageUp
	KeyPageDown

	KeyCtrlEnter
	KeyCtrlBackspace

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

var KeyToCanonical = map[Key]Key{
	KeyCtrlH: KeyCtrlBackspace,
	KeyCtrlI: KeyTab,
	KeyCtrlJ: KeyCtrlEnter,
	KeyCtrlM: KeyEnter,
	KeyCtrlZ: KeyPause,
}

func (k Key) Normalize() Key {
	if canonical, ok := KeyToCanonical[k]; ok {
		return canonical
	}

	return k
}

func (k Key) Equivalents() []Key {
	results := []Key{k}

	if canonical, ok := KeyToCanonical[k]; ok {
		results = append(results, canonical)
	}

	for combo, canonical := range KeyToCanonical {
		if canonical == k {
			results = append(results, combo)
		}
	}

	return results
}

var SequenceToKey = map[string]Key{
	"": KeyUnknown,

	"\t":      KeyTab,
	CSI + "Z": KeyShiftTab,

	CSI + "A": KeyUp,
	CSI + "B": KeyDown,
	CSI + "C": KeyRight,
	CSI + "D": KeyLeft,

	CSI + "H": KeyHome,
	CSI + "F": KeyEnd,

	CSI + "1;5A": KeyCtrlUp,
	CSI + "1;5B": KeyCtrlDown,
	CSI + "1;5C": KeyCtrlRight,
	CSI + "1;5D": KeyCtrlLeft,

	CSI + "1;3A": KeyAltUp,
	CSI + "1;3B": KeyAltDown,
	CSI + "1;3C": KeyAltRight,
	CSI + "1;3D": KeyAltLeft,

	CSI + "1;2A": KeyShiftUp,
	CSI + "1;2B": KeyShiftDown,
	CSI + "1;2C": KeyShiftRight,
	CSI + "1;2D": KeyShiftLeft,

	"\x0d":     KeyEnter,
	"\x7f":     KeyBackspace,
	"\x1a":     KeyPause,
	ESC:        KeyEsc,
	CSI + "2~": KeyInsert,
	CSI + "3~": KeyDelete,
	CSI + "5~": KeyPageUp,
	CSI + "6~": KeyPageDown,

	"\x0a": KeyCtrlEnter,
	"\x08": KeyCtrlBackspace,

	ESC + "OP":  KeyF1,
	ESC + "OQ":  KeyF2,
	ESC + "OR":  KeyF3,
	ESC + "OS":  KeyF4,
	CSI + "15~": KeyF5,
	CSI + "17~": KeyF6,
	CSI + "18~": KeyF7,
	CSI + "19~": KeyF8,
	CSI + "20~": KeyF9,
	CSI + "21~": KeyF10,
	CSI + "23~": KeyF11,
	CSI + "24~": KeyF12,

	"\x01": KeyCtrlA,
	"\x02": KeyCtrlB,
	"\x03": KeyCtrlC,
	"\x04": KeyCtrlD,
	"\x05": KeyCtrlE,
	"\x06": KeyCtrlF,
	"\x07": KeyCtrlG,
	//"\x08": KeyCtrlH,
	//"\x09": KeyCtrlI,
	//"\x0a": KeyCtrlJ,
	"\x0b": KeyCtrlK,
	"\x0c": KeyCtrlL,
	//"\x0d": KeyCtrlM,
	"\x0e": KeyCtrlN,
	"\x0f": KeyCtrlO,
	"\x10": KeyCtrlP,
	"\x11": KeyCtrlQ,
	"\x12": KeyCtrlR,
	"\x13": KeyCtrlS,
	"\x14": KeyCtrlT,
	"\x15": KeyCtrlU,
	"\x16": KeyCtrlV,
	"\x17": KeyCtrlW,
	"\x18": KeyCtrlX,
	"\x19": KeyCtrlY,
	//"\x1a": KeyCtrlZ,

	ESC + "a": KeyAltA,
	ESC + "b": KeyAltB,
	ESC + "c": KeyAltC,
	ESC + "d": KeyAltD,
	ESC + "e": KeyAltE,
	ESC + "f": KeyAltF,
	ESC + "g": KeyAltG,
	ESC + "h": KeyAltH,
	ESC + "i": KeyAltI,
	ESC + "j": KeyAltJ,
	ESC + "k": KeyAltK,
	ESC + "l": KeyAltL,
	ESC + "m": KeyAltM,
	ESC + "n": KeyAltN,
	ESC + "o": KeyAltO,
	ESC + "p": KeyAltP,
	ESC + "q": KeyAltQ,
	ESC + "r": KeyAltR,
	ESC + "s": KeyAltS,
	ESC + "t": KeyAltT,
	ESC + "u": KeyAltU,
	ESC + "v": KeyAltV,
	ESC + "w": KeyAltW,
	ESC + "x": KeyAltX,
	ESC + "y": KeyAltY,
	ESC + "z": KeyAltZ,
}
