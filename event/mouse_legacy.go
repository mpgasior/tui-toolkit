package event

import (
	"bytes"

	"github.com/mpgasior/tui-toolkit/ansi"
)

func (d *MouseDecoder) decodeMouseLegacy(buf []byte) (mouse Mouse, ok bool) {
	if !bytes.HasPrefix(buf, ansi.MouseLegacyPrefixBytes) {
		return mouse, false
	}

	if len(buf) != ansi.MouseLegacyLength {
		return mouse, false
	}

	x := buf[ansi.MouseLegacyColIndex] - ansi.MouseLegacyOffset
	y := buf[ansi.MouseLegacyRowIndex] - ansi.MouseLegacyOffset
	ctrl := buf[ansi.MouseLegacyCtrlIndex] - ansi.MouseLegacyOffset

	var (
		button MouseButton
		action MouseAction
		mods   MouseModifier
	)

	const (
		maskButton byte = 0x3
		flagShift  byte = 1 << 2
		flagAlt    byte = 1 << 3
		flagCtrl   byte = 1 << 4
		flagMotion byte = 1 << 5
		flagWheel  byte = 1 << 6
	)

	if ctrl&flagShift != 0 {
		mods |= MouseModifierShift
	}

	if ctrl&flagAlt != 0 {
		mods |= MouseModifierAlt
	}

	if ctrl&flagCtrl != 0 {
		mods |= MouseModifierCtrl
	}

	rawBtn := ctrl & maskButton
	isWheel := ctrl&flagWheel != 0
	isMotion := ctrl&flagMotion != 0

	switch {
	case isWheel:
		button = MouseButtonWheelUp
		if rawBtn != 0 {
			button = MouseButtonWheelDown
		}
	case rawBtn == 3:
		button = MouseButtonNone
		if isMotion {
			action = MouseActionMove
		} else {
			action = MouseActionRelease
		}
	default:
		button = MouseButton(rawBtn + 1)
	}

	mouse = Mouse{
		X: int(x), Y: int(y),
		Button: button,
		Action: action,
		Mods:   mods,
	}

	return mouse, true
}
