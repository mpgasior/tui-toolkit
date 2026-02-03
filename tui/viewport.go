package tui

import (
	"sort"
)

type ViewportConstraint struct {
	name     string
	isStatic bool
	value    int
}

func ViewportConstraintFixed(name string, value int) ViewportConstraint {
	return ViewportConstraint{
		name: name, isStatic: true, value: value,
	}
}

func ViewportConstraintDynamic(name string, value int) ViewportConstraint {
	return ViewportConstraint{
		name: name, isStatic: false, value: value,
	}
}

type Viewport struct {
	buffer *Buffer
	x, y   int
	w, h   int
}

func NewViewport(buf *Buffer) Viewport {
	w, h := buf.Size()

	return Viewport{
		buffer: buf,
		x:      0, y: 0,
		w: w, h: h,
	}
}

func (vp Viewport) Size() (w, h int) {
	return vp.w, vp.h
}

func (vp Viewport) SetAt(x, y int, primary rune, combs []rune, width uint8, style Style) {
	if x < 0 || x >= vp.w || y < 0 || y >= vp.h {
		return
	}

	vp.buffer.SetAt(x+vp.x, y+vp.y, primary, combs, width, style)
}

func (vp Viewport) SetCursorPos(x, y int) {
	if x < 0 || x >= vp.w || y < 0 || y >= vp.h {
		vp.buffer.SetCursorPos(-1, -1)
		return
	}
	vp.buffer.SetCursorPos(vp.x+x, vp.y+y)
}

func (vp Viewport) Slice(x, y, w, h int) Viewport {
	originX := max(0, min(x, vp.w))
	originY := max(0, min(y, vp.h))

	remainingW := vp.w - originX
	remainingH := vp.h - originY

	return Viewport{
		buffer: vp.buffer,
		x:      vp.x + originX,
		y:      vp.y + originY,
		w:      max(0, min(w, remainingW)),
		h:      max(0, min(h, remainingH)),
	}
}

func (vp Viewport) Offset(offsets ...int) Viewport {
	var top, right, bottom, left int

	switch len(offsets) {
	case 1:
		top, right, bottom, left = offsets[0], offsets[0], offsets[0], offsets[0]
	case 2:
		top, bottom = offsets[0], offsets[0]
		right, left = offsets[1], offsets[1]
	case 3:
		top = offsets[0]
		right, left = offsets[1], offsets[1]
		bottom = offsets[2]
	case 4:
		top, right, bottom, left = offsets[0], offsets[1], offsets[2], offsets[3]
	}

	newX := vp.x + left
	newY := vp.y + top

	newW := max(0, vp.w-left-right)
	newH := max(0, vp.h-top-bottom)

	return Viewport{
		buffer: vp.buffer,
		x:      newX,
		y:      newY,
		w:      newW,
		h:      newH,
	}
}

func distribute(total int, weights ...int) []int {
	type allocation struct {
		index    int
		value    int
		fraction int
	}

	sumWeights := 0
	for _, w := range weights {
		sumWeights += w
	}

	allocations := make([]allocation, len(weights))
	distributed := 0

	for i, w := range weights {
		value := (total * w) / sumWeights
		fraction := (total * w) % sumWeights

		allocations[i] = allocation{
			index:    i,
			value:    value,
			fraction: fraction,
		}

		distributed += value
	}

	leftover := total - distributed

	sort.Slice(allocations, func(i, j int) bool {
		return allocations[i].fraction > allocations[j].fraction
	})

	for i := range leftover {
		allocations[i].value += 1
	}

	result := make([]int, len(weights))
	for _, allocation := range allocations {
		result[allocation.index] = allocation.value
	}
	return result
}

func (vp Viewport) splitMap(horizontal bool, cs ...ViewportConstraint) map[string]Viewport {
	fixedTotal := 0
	weights := make([]int, len(cs))
	for i, c := range cs {
		if c.isStatic {
			fixedTotal += c.value
		} else {
			weights[i] = c.value
		}
	}

	total := vp.w
	if horizontal {
		total = vp.h
	}

	allocations := distribute(total-fixedTotal, weights...)
	result := make(map[string]Viewport)

	offset := 0
	for i, allocation := range allocations {
		c := cs[i]
		if c.name == "" {
			continue
		}

		if allocation == 0 {
			allocation = c.value
		}

		x, y, w, h := offset, 0, allocation, vp.h
		if horizontal {
			x, y, w, h = 0, offset, vp.w, allocation
		}

		result[c.name] = vp.Slice(x, y, w, h)
		offset += allocation
	}

	return result
}

func (vp Viewport) SplitH(cs ...ViewportConstraint) map[string]Viewport {
	return vp.splitMap(true, cs...)
}

func (vp Viewport) SplitV(cs ...ViewportConstraint) map[string]Viewport {
	return vp.splitMap(false, cs...)
}
