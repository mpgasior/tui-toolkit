package termio

import (
	"io"

	"github.com/mpgasior/tui-toolkit/ansi"
)

var (
	WithAltScreen      = ModBytes(ansi.AltScreenOnBytes, ansi.AltScreenOffBytes)
	WithBracketedPaste = ModBytes(ansi.BracketedPasteOnBytes, ansi.BracketedPasteOffBytes)
	WithCursorHidden   = ModBytes(ansi.CursorHideBytes, ansi.CursorShowBytes)
	WithMouseNormal    = ModBytes(ansi.MouseNormalModeOnBytes, ansi.MouseNormalModeOffBytes)
	WithMouseAny       = ModBytes(ansi.MouseAnyModeOnBytes, ansi.MouseAnyModeOffBytes)
	WithMouseSGR       = ModBytes(ansi.MouseSGRModeOnBytes, ansi.MouseSGRModeOffBytes)
)

func ModString(apply, restore string) SessionModifier {
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

func ModBytes(apply, restore []byte) SessionModifier {
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
