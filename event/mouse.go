package event

type MouseButton uint8

const (
	MouseButtonNone MouseButton = iota
	MouseButtonLeft
	MouseButtonMiddle
	MouseButtonRight
	MouseButtonWheelUp
	MouseButtonWheelDown
)

type MouseAction uint8

const (
	MouseActionPress MouseAction = iota
	MouseActionRelease
	MouseActionDrag
	MouseActionMove
)

type MouseModifier uint8

const (
	MouseModifierShift = 1 << iota
	MouseModifierAlt
	MouseModifierCtrl
)

type Mouse struct {
	X, Y   int
	Button MouseButton
	Action MouseAction
	Mods   MouseModifier
}

type MouseDecoder struct {
}

func NewMouseDecoder() *MouseDecoder {
	return &MouseDecoder{}
}

func (d *MouseDecoder) Decode(buf []byte) (mouse Mouse, ok bool) {
	if legacy, ok := d.decodeMouseLegacy(buf); ok {
		return legacy, ok
	}

	if modern, ok := d.decodeMouseSGR(buf); ok {
		return modern, ok
	}

	return mouse, true
}
