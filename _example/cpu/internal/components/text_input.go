package components

import (
	"unicode/utf8"

	"github.com/mpgasior/tui-toolkit/vt"
)

type TextInput struct {
	buffer []rune
	cursor int
}

func (t *TextInput) Update(key vt.KeyEvent) (consumed bool) {
	switch key.Key {
	case vt.KeyBackspace:
		if t.cursor > 0 {
			t.buffer = append(t.buffer[:t.cursor-1], t.buffer[t.cursor:]...)
			return true
		}
	case vt.KeyLeft:
		if t.cursor > 0 {
			t.cursor -= 1
			return true
		}
	case vt.KeyRight:
		if t.cursor < len(t.buffer) {
			t.cursor += 1
			return true
		}
	default:
		if key.Rune != 0 {
			t.buffer = append(t.buffer[:t.cursor], key.Rune, t.buffer[t.cursor:]...)
			t.cursor += 1
			return true
		}
	}

	return false
}

func (t *TextInput) String() string {
	return string(t.buffer)
}

func (t *TextInput) Clear() {
	t.cursor = 0
	t.buffer = nil
}
