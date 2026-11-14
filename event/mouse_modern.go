package event

import (
	"bytes"
	"strconv"

	"github.com/mpgasior/tui-toolkit/ansi"
)

func (d *MouseDecoder) decodeMouseSGR(buf []byte) (mouse Mouse, ok bool) {
	if !bytes.HasPrefix(buf, ansi.MouseSGRPrefixBytes) {
		return mouse, false
	}

	if !ansi.MouseSGRFinalRange.Contains(buf[len(buf)-1]) {
		return mouse, false
	}

	isRelease := buf[len(buf)-1] == 'm'

	start, end := len(ansi.MouseSGRPrefixBytes), len(buf)-1
	parts := bytes.Split(buf[start:end], []byte{';'})
	if len(parts) != 3 {
		return mouse, false
	}

	ctrl, err := strconv.Atoi(string(parts[0]))
	if err != nil {
		return mouse, false
	}

	x, err := strconv.Atoi(string(parts[1]))
	if err != nil {
		return mouse, false
	}

	y, err := strconv.Atoi(string(parts[2]))
	if err != nil {
		return mouse, false
	}

	var (
		button MouseButton
		action MouseAction
		mods   MouseModifier
	)

	const (
		maskButton = 0x3
		flagShift  = 1 << 2
		flagAlt    = 1 << 3
		flagCtrl   = 1 << 4
		flagMotion = 1 << 5
		flagWheel  = 1 << 6
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
		action = MouseActionPress
		button = MouseButtonWheelUp
		if rawBtn != 0 {
			button = MouseButtonWheelDown
		}
	case rawBtn == 3:
		action = MouseActionMove
		button = MouseButtonNone
	default:
		button = MouseButton(rawBtn + 1)
		action = MouseActionPress
		switch {
		case isRelease:
			action = MouseActionRelease
		case isMotion:
			action = MouseActionDrag
		}
	}

	mouse = Mouse{
		X: x, Y: y,
		Button: button,
		Action: action,
		Mods:   mods,
	}

	return mouse, true
}
