package draw

import (
	"bufio"
	"strings"

	"github.com/mpgasior/tui-toolkit/screen"
)

func Text(b screen.Buffer, text string, style screen.Style) {
	scanner := bufio.NewScanner(strings.NewReader(text))

	row := 0
	for scanner.Scan() {
		line := scanner.Text()

		for idx, r := range []rune(line) {
			Rune(b, idx, row, r, style)
		}
		row += 1
	}
}
