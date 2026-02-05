package draw

import (
	"bufio"
	"strings"

	"github.com/mpgasior/tui-toolkit/screen"
)

func Text(m screen.Buffer, text string, style screen.Style) {
	scanner := bufio.NewScanner(strings.NewReader(text))

	row := 0
	for scanner.Scan() {
		line := scanner.Text()

		for idx, r := range []rune(line) {
			Rune(m, idx, row, r, style)
		}
		row += 1
	}
}
