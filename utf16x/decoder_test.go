package utf16x_test

import (
	"testing"

	"github.com/mpgasior/tui-toolkit/utf16x"
)

func TestDecoder(t *testing.T) {
	tests := []struct {
		name   string
		input  []uint16
		want   rune
		wantOk bool
	}{
		{
			name:   "Standard BMP character (Z)",
			input:  []uint16{0x005A},
			want:   'Z',
			wantOk: true,
		},
		{
			name:   "Valid Surrogate Pair (Grinning Face Emoji)",
			input:  []uint16{0xD83D, 0xDE00}, // 😊
			want:   '😀',
			wantOk: true,
		},
		{
			name:   "First half of surrogate only",
			input:  []uint16{0xD83D},
			want:   0,
			wantOk: false,
		},
		{
			name:   "High Surrogate followed by character 'Z'",
			input:  []uint16{0xD83D, 0x005A},
			want:   'Z',
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := utf16x.Decoder{}
			var got rune
			var ok bool

			for _, val := range tt.input {
				got, ok = decoder.Decode(val)
			}

			if got != tt.want || ok != tt.wantOk {
				t.Errorf("Decode() = (%U, %v), want (%U, %v)", got, ok, tt.want, tt.wantOk)
			}
		})
	}
}
