package utf16x

import (
	"testing"
)

func TestRuneScanner(t *testing.T) {
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
			rs := &RuneScanner{}
			var got rune
			var ok bool

			for _, val := range tt.input {
				got, ok = rs.Scan(val)
			}

			if got != tt.want || ok != tt.wantOk {
				t.Errorf("Write() = (%U, %v), want (%U, %v)", got, ok, tt.want, tt.wantOk)
			}
		})
	}
}
