package termio

import (
	"io"

	"github.com/mpgasior/tui-toolkit/ansi"
)

var (
	WithAltScreen      = Mod(ansi.AltScreenOn, ansi.AltScreenOff)
	WithBracketedPaste = Mod(ansi.BracketedPasteOn, ansi.BracketedPasteOff)
	WithCursorHidden   = Mod(ansi.CursorHide, ansi.CursorShow)
)

func Mod(apply, restore string) SessionModifier {
	return SessionModifier{
		Apply: func(s Session) error {
			_, err := io.WriteString(s, apply)
			return err
		},
		Restore: func(s Session) error {
			_, err := io.WriteString(s, restore)
			return err
		},
	}
}
