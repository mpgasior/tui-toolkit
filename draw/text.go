package draw

import (
	"bufio"
	"strings"
	"unicode/utf8"

	"github.com/mpgasior/tui-toolkit/screen"
)

type TextAlignment int

const (
	TextAlignmentLeft TextAlignment = iota
	TextAlignmentCenter
	TextAlignmentRight
)

type TextChunk struct {
	Text      string
	Style     screen.Style
	Alignment TextAlignment
}

func Text(m screen.Mutator, chunks ...TextChunk) {
	if len(chunks) == 0 {
		return
	}

	screenWidth, _ := m.Size()

	textWidth := 0
	for _, chunk := range chunks {
		textWidth += utf8.RuneCountInString(chunk.Text)
	}
	textWidth = min(textWidth, screenWidth)

	offset := 0
	switch chunks[0].Alignment {
	case TextAlignmentCenter:
		offset = (screenWidth - textWidth) / 2
	case TextAlignmentRight:
		offset = screenWidth - textWidth
	case TextAlignmentLeft:
		fallthrough
	default:
		offset = 0
	}

	for _, chunk := range chunks {
		localOffset := 0
		for idx, r := range chunk.Text {
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

		for idx, r := range line {
			Rune(m, idx, row, r, style)
		}

		row += 1
	}
}

func Line(m screen.Mutator, line string, style screen.Style) {
	for idx, r := range line {
		Rune(m, idx, 0, r, style)
	}
}
