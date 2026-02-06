package draw

import (
	"bufio"
	"strings"

	"github.com/mpgasior/tui-toolkit/screen"
)

type TextChunk struct {
	Text  string
	Style screen.Style
}

func Text(m screen.Mutator, chunks ...TextChunk) {
	offset := 0
	for _, chunk := range chunks {
		localOffset := 0
		for idx, r := range []rune(chunk.Text) {
			Rune(m, offset+idx, 0, r, chunk.Style)
			localOffset += 1
		}
		offset += localOffset
	}
}

func Lines(m screen.Mutator, text string, style screen.Style) {
	scanner := bufio.NewScanner(strings.NewReader(text))

	row := 0
	for scanner.Scan() {
		line := scanner.Text()

		Line(m, line, style)

		row += 1
	}
}

func Line(m screen.Mutator, line string, style screen.Style) {
	for idx, r := range []rune(line) {
		Rune(m, idx, 0, r, style)
	}
}
