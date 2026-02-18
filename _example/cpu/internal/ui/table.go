package ui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/model"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
	"github.com/mpgasior/tui-toolkit/vt"
)

var tableColumnOrder = []model.SortBy{
	model.SortByPID,
	model.SortByAvgCPU,
	model.SortByRecentCPU,
	model.SortByAge,
	model.SortByName,
}

type Table struct {
	Rows      []model.QueryResult
	SortBy    model.SortBy
	SortOrder model.SortOrder
	IsPaused  bool
	Scroll    view.Scroll
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

func (t *Table) ResetBusy() {
	t.IsPaused = false
	t.Scroll.Jump(0)
}

func (t *Table) Draw(vp view.Port, focused bool) {
	boxStyle := screen.DefaultStyle
	if focused {
		boxStyle = boxStyle.Fg(screen.ColorGreen)
	}
	draw.Box(vp, draw.BoxBorderThin, boxStyle)
	draw.Clear(vp.Offset(0, 2, 1, 2), boxStyle)

	layout := view.SplitV(vp,
		view.Fixed("selected", 4),
		view.Fixed("pid", 7),
		view.Fixed("avg-cpu", 13),
		view.Fixed("", 2),
		view.Fixed("recent-cpu", 10),
		view.Fixed("", 2),
		view.Fixed("memory", 10),
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
	drawHeader("recent-cpu", "CPU% (Now)", model.SortByRecentCPU)
	drawHeader("age", "Age", model.SortByAge)

	_, h := vp.Size()
	t.drawFooter(focused, cell, h, boxStyle)

	if t.Rows == nil {
		return
	}

	start, end := t.Scroll.Update(h-2, len(t.Rows))
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
		draw.Text(cell("age", rowIdx), draw.TextChunk{
			Text:  formatAge(row.Age),
			Style: rowStyle,
		})
		draw.Text(cell("avg-cpu", rowIdx), draw.TextChunk{
			Text:      fmt.Sprintf("%5.2f%%", row.AvgCPU),
			Style:     rowStyle,
			Alignment: draw.TextAlignmentRight,
		})
		draw.Text(cell("recent-cpu", rowIdx), draw.TextChunk{
			Text:      fmt.Sprintf("%5.2f%%", row.RecentCPU),
			Style:     rowStyle,
			Alignment: draw.TextAlignmentRight,
		})

		rowIdx += 1
	}
}

func (t *Table) drawFooter(focused bool, cell func(key string, row int) view.Port, h int, boxStyle screen.Style) {
	text := "Total: " + strconv.FormatInt(int64(len(t.Rows)), 10)
	if focused {
		text = strconv.FormatInt(int64(t.Scroll.Index+1), 10) + " of " + strconv.FormatInt(int64(len(t.Rows)), 10)
	}

	if t.IsPaused {
		text = "[Paused] " + text
	}

	draw.Text(cell("name", h-1).Offset(0, 2), draw.TextChunk{
		Text:      text,
		Style:     boxStyle,
		Alignment: draw.TextAlignmentRight,
	})
}

func formatAge(d time.Duration) string {
	if d == 0 {
		return "N/A"
	}

	if d.Hours() >= 24 {
		days := int(d.Hours() / 24)
		hours := int(d.Hours()) % 24
		return fmt.Sprintf("%dd %dh", days, hours)
	}

	return d.Round(time.Second).String()
}
