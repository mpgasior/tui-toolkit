package ui

import (
	"fmt"
	"strconv"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/model"
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
	"github.com/mpgasior/tui-toolkit/vt"
)

var tableColumnOrder = []model.SortBy{
	model.SortByPID,
	model.SortByAvgCPU,
	model.SortByCPU,
	model.SortByMem,
	model.SortByMaxMem,
	model.SortByAge,
	model.SortByName,
}

type Table struct {
	Rows      []model.Process
	SortBy    model.SortBy
	SortOrder model.SortOrder

	Selected process.Key
	Scroll   view.Scroll
}

func NewTable() Table {
	return Table{
		Scroll: view.Scroll{
			Margin: 2,
		},
	}
}

func (t *Table) Update(key vt.KeyEvent) (didUpdate bool) {
	switch key.Rune {
	case 's':
		t.NextSortOrder()
		return true
	case 'h':
		t.PrevSortBy()
		return true
	case 'l':
		t.NextSortBy()
		return true
	}

	switch key.Key {
	case vt.KeyEnter:
		count := len(t.Rows)
		if count > 0 && t.Scroll.Index < count {
			t.Selected = t.Rows[t.Scroll.Index].Key
			return true
		}
	case vt.KeyJ:
		t.Scroll.Move(1)
	case vt.KeyK:
		t.Scroll.Move(-1)
	case vt.KeyG:
		t.Scroll.Jump(0)
	case vt.KeyShiftG:
		t.Scroll.Jump(len(t.Rows) - 1)
	case vt.KeyCtrlU:
		t.Scroll.Move(-5)
	case vt.KeyCtrlD:
		t.Scroll.Move(5)
	}

	return false
}

func (t *Table) NextSortOrder() {
	t.SortOrder = (t.SortOrder + 1) % 2
}

func (t *Table) NextSortBy() {
	for i, sort := range tableColumnOrder {
		if sort == t.SortBy {
			t.SortBy = tableColumnOrder[(i+1)%len(tableColumnOrder)]
			break
		}
	}
}

func (t *Table) PrevSortBy() {
	n := len(tableColumnOrder)
	for i, sort := range tableColumnOrder {
		if sort == t.SortBy {
			t.SortBy = tableColumnOrder[(i-1+n)%n]
			break
		}
	}
}

func (t *Table) Reset() {
	t.Selected = process.KeyNone
	t.Scroll.Jump(0)
}

func (t *Table) Draw(vp view.Port, focused, paused bool) {
	activeStyle := screen.DefaultStyle
	if focused {
		activeStyle = activeStyle.Fg(screen.ColorGreen)
	}
	draw.Box(vp, draw.BoxBorderThin, activeStyle)
	draw.Clear(vp.Offset(0, 2, 1, 2), activeStyle)

	layout := view.SplitV(vp,
		view.Fixed("selected", 4),
		view.Fixed("pid", 7),
		view.Fixed("avg-cpu", 13),
		view.Fixed("", 2),
		view.Fixed("recent-cpu", 10),
		view.Fixed("", 2),
		view.Fixed("recent-mem", 10),
		view.Fixed("", 2),
		view.Fixed("peak-mem", 10),
		view.Fixed("", 2),
		view.Fixed("age", 10),
		view.Fixed("", 2),
		view.Dynamic("name", 15),
	)

	cell := func(key string, row int) view.Port {
		col := layout[key]
		w, _ := col.Size()
		return col.Slice(0, row, w, 1)
	}

	drawHeader := func(key, label string, sortBy model.SortBy) {
		style := screen.DefaultStyle.
			Attr(screen.AttrBold | screen.AttrUnderline)

		if sortBy == t.SortBy {
			style = style.Fg(screen.ColorGreen)
			r := '↑'
			if t.SortOrder == model.SortOrderDescending {
				r = '↓'
			}

			draw.Line(cell(key, 0), string(r)+label, style)
			return
		}

		draw.Text(cell(key, 0), draw.TextChunk{
			Text:  label,
			Style: style,
		})
	}

	drawHeader("pid", "PID", model.SortByPID)
	drawHeader("name", "Name", model.SortByName)
	drawHeader("avg-cpu", "CPU% (Avg 1m)", model.SortByAvgCPU)
	drawHeader("recent-cpu", "CPU% (Now)", model.SortByCPU)
	drawHeader("recent-mem", "MEM (Now)", model.SortByMem)
	drawHeader("peak-mem", "MEM (Peak)", model.SortByMaxMem)
	drawHeader("age", "Age", model.SortByAge)

	_, h := vp.Size()
	start, end := t.Scroll.Update(h-2, len(t.Rows))

	t.drawFooter(cell("name", h-1), focused, paused, activeStyle)
	t.drawScroll(vp, activeStyle)

	if t.Rows == nil {
		return
	}

	rows := t.Rows[start:end]

	rowIdx := 1
	for idx, row := range rows {
		rowStyle := screen.DefaultStyle
		if idx+t.Scroll.Offset == t.Scroll.Index {
			rowStyle = rowStyle.Fg(screen.ColorCyan)
			draw.Text(cell("selected", rowIdx), draw.TextChunk{
				Text:      "┃",
				Style:     screen.DefaultStyle.Fg(screen.ColorRed),
				Alignment: draw.TextAlignmentCenter,
			})
		}

		draw.Text(cell("pid", rowIdx), draw.TextChunk{
			Text:  strconv.FormatInt(int64(row.PID), 10),
			Style: rowStyle,
		})
		draw.Text(cell("name", rowIdx), draw.TextChunk{
			Text:  row.Name,
			Style: rowStyle,
		})

		if row.AgeReady {
			text := formatDuration(row.Age)
			if !row.ExitTime.IsZero() {
				text = "dead"
			}

			draw.Text(cell("age", rowIdx), draw.TextChunk{
				Text:  text,
				Style: rowStyle,
			})
		}

		if row.MemReady {
			draw.Text(cell("peak-mem", rowIdx), draw.TextChunk{
				Text:      formatWorkingSet(row.MemoryMax),
				Style:     rowStyle,
				Alignment: draw.TextAlignmentRight,
			})
			draw.Text(cell("recent-mem", rowIdx), draw.TextChunk{
				Text:      formatWorkingSet(row.MemoryRSS),
				Style:     rowStyle,
				Alignment: draw.TextAlignmentRight,
			})
		}

		if row.CPUReady {
			draw.Text(cell("avg-cpu", rowIdx), draw.TextChunk{
				Text:      fmt.Sprintf("%.2f%%", row.CPUAvg),
				Style:     rowStyle,
				Alignment: draw.TextAlignmentRight,
			})
			draw.Text(cell("recent-cpu", rowIdx), draw.TextChunk{
				Text:      fmt.Sprintf("%.2f%%", row.CPU),
				Style:     rowStyle,
				Alignment: draw.TextAlignmentRight,
			})
		}

		rowIdx += 1
	}
}

func (t *Table) drawScroll(vp view.Port, style screen.Style) {
	w, h := vp.Size()
	totalRows := len(t.Rows)
	if totalRows == 0 {
		return
	}

	trackHeight := h - 2
	if trackHeight <= 0 {
		return
	}

	thumbHeight := int(float64(trackHeight) * float64(trackHeight) / float64(totalRows))
	if thumbHeight == 0 {
		thumbHeight = 1
	}

	maxOffset := totalRows - trackHeight
	scrollRatio := float64(t.Scroll.Offset) / float64(maxOffset)

	startPos := int(scrollRatio * float64(trackHeight-thumbHeight))

	scrollBar := vp.Slice(w-1, 0, 1, h).Offset(1, 0, 1, 0)
	for idx := 0; idx < thumbHeight; idx += 1 {
		draw.Rune(scrollBar, 0, startPos+idx, '▐', style)
	}
}

func (t *Table) drawFooter(vp view.Port, focused, paused bool, style screen.Style) {
	text := "Total: " + strconv.FormatInt(int64(len(t.Rows)), 10)
	if focused {
		index := strconv.FormatInt(int64(t.Scroll.Index+1), 10)
		total := strconv.FormatInt(int64(len(t.Rows)), 10)
		text = index + " of " + total
	}

	if paused {
		text = "[Paused] " + text
	}

	draw.Text(vp.Offset(0, 2), draw.TextChunk{
		Text:      text,
		Style:     style,
		Alignment: draw.TextAlignmentRight,
	})
}
