package termio

import "github.com/mpgasior/tui-toolkit/ansi"

var (
	WithAltScreen      = Mod(ansi.AltScreenOn, ansi.AltScreenOff)
	WithBracketedPaste = Mod(ansi.BracketedPasteOn, ansi.BracketedPasteOff)
	WithCursorHidden   = Mod(ansi.CursorHide, ansi.CursorShow)
)

func Mod(apply []byte, restore []byte) SessionModifier {
	return SessionModifier{
		Apply: func(s Session) error {
			_, err := s.Write(apply)
			return err
		},
		Restore: func(s Session) error {
			_, err := s.Write(restore)
			return err
		},
	}
}
