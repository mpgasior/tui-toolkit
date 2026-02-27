package draw

import (
	"math"

	"github.com/mpgasior/tui-toolkit/screen"
)

var HistogramBlocks = [9]rune{' ', '▂', '▃', '▄', '▅', '▆', '▇', '█', '█'}

func Histogram[T any](m screen.Mutator, data []T, fn func(t T) float64, style screen.Style) {
	if len(data) == 0 {
		return
	}

	w, h := m.Size()

	first := fn(data[0])
	dataMin, dataMax := first, first

	for _, v := range data {
		valF := fn(v)
		if valF < dataMin {
			dataMin = valF
		}
		if valF > dataMax {
			dataMax = valF
		}
	}

	dataRange := dataMax - dataMin
	rangeMax := float64(h)

	for idx, item := range data {
		if idx >= w {
			break
		}

		var v float64
		if dataRange == 0 {
			v = float64(h) * 0.5
		} else {
			v = ((fn(item) - dataMin) / dataRange) * rangeMax
		}

		fullBlocks := int(math.Floor(v))
		remainder := v - float64(fullBlocks)
		partialIdx := int(remainder * 8)

		for yOffset := 0; yOffset < fullBlocks && yOffset < h-1; yOffset++ {
			Rune(m, idx, h-1-yOffset, '█', style)
		}

		if fullBlocks < h {
			char := HistogramBlocks[partialIdx]
			if char != ' ' {
				Rune(m, idx, h-1-fullBlocks, char, style)
			}
		}
	}
}
