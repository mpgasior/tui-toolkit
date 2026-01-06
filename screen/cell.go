package screen

import "slices"

type Cell struct {
	Primary   rune
	Combining []rune
	Style     Style
	Width     uint8
}

func (c *Cell) Equal(other *Cell) bool {
	if c.Primary != other.Primary {
		return false
	}
	if c.Style != other.Style {
		return false
	}
	if c.Width != other.Width {
		return false
	}
	if len(c.Combining) != len(other.Combining) {
		return false
	}

	return slices.Equal(c.Combining, other.Combining)
}
