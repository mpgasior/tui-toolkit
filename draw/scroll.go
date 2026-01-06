package draw

import "github.com/mpgasior/tui-toolkit/screen"

var RuneScrollVRight rune = '▐'

func ScrollV(m screen.Mutator, offset, total int, r rune, style screen.Style) {
	if total == 0 {
		return
	}

	_, h := m.Size()
	if h <= 0 {
		return
	}

	thumbHeight := int(float64(h) * float64(h) / float64(total))
	if thumbHeight == 0 {
		thumbHeight = 1
	}

	maxOffset := total - h
	scrollRatio := float64(offset) / float64(maxOffset)

	startPos := int(scrollRatio * float64(h-thumbHeight))

	for idx := 0; idx < thumbHeight; idx += 1 {
		Rune(m, 0, startPos+idx, r, style)
	}
}
