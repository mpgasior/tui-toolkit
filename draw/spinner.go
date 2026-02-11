package draw

import "github.com/mpgasior/tui-toolkit/screen"

type SpinnerFrame []rune

var (
	SpinnerBraille = SpinnerFrame{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
	SpinnerLine    = SpinnerFrame{'|', '/', '-', '\\'}
	SpinnerBlocks  = SpinnerFrame{'▖', '▘', '▝', '▗'}
)

func Spinner(m screen.Mutator, i int, frame SpinnerFrame, style screen.Style) {
	r := frame[i%len(frame)]

	Rune(m, 0, 0, r, style)
}
