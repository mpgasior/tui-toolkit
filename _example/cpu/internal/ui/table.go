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

var TableColumnOrder = []model.SortBy{
	model.SortByPID,
	model.SortByAvgCPU,
	model.SortByRecentCPU,
	model.SortByCreationTime,
	model.SortByName,
}

func NextSortBy(current model.SortBy) model.SortBy {
	for i, sort := range TableColumnOrder {
		if sort == current {
			return TableColumnOrder[(i+1)%len(TableColumnOrder)]
		}
	}
	return TableColumnOrder[0]
}

func PrevSortBy(current model.SortBy) model.SortBy {
	n := len(TableColumnOrder)
	for i, sort := range TableColumnOrder {
		if sort == current {
			return TableColumnOrder[(i-1+n)%n]
		}
	}
	return TableColumnOrder[n-1]
}

type Table struct {
	Rows      []model.QueryResult
	SortBy    model.SortBy
	SortOrder model.SortOrder
}

func (t *Table) Update(key vt.KeyEvent) (didUpdate bool) {
	switch key.Rune {
	case 's':
		t.SortOrder = (t.SortOrder + 1) % 2
		return true
	case 'h':
		t.SortBy = PrevSortBy(t.SortBy)
		return true
	case 'l':
		t.SortBy = NextSortBy(t.SortBy)
		return true
	}

	return false
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
	drawHeader("age", "Age", model.SortByCreationTime)

	_, h := vp.Size()

	for idx, row := range t.Rows {
		if idx >= h-2 {
			break
		}
		draw.Text(cell("pid", idx+1), draw.TextChunk{
			Text:  strconv.FormatInt(int64(row.PID), 10),
			Style: screen.DefaultStyle,
		})
		draw.Text(cell("name", idx+1), draw.TextChunk{
			Text:  row.Name,
			Style: screen.DefaultStyle,
		})
		draw.Text(cell("age", idx+1), draw.TextChunk{
			Text:  formatCompact(row.CreationTime),
			Style: screen.DefaultStyle,
		})
		draw.Text(cell("avg-cpu", idx+1), draw.TextChunk{
			Text:      fmt.Sprintf("%5.2f%%", row.AvgCPU),
			Style:     screen.DefaultStyle,
			Alignment: draw.TextAlignmentRight,
		})
		draw.Text(cell("recent-cpu", idx+1), draw.TextChunk{
			Text:      fmt.Sprintf("%5.2f%%", row.RecentCPU),
			Style:     screen.DefaultStyle,
			Alignment: draw.TextAlignmentRight,
		})
	}

	draw.Text(cell("name", h-1).Offset(0, 2), draw.TextChunk{
		Text:      "Total: " + strconv.FormatInt(int64(len(t.Rows)), 10),
		Style:     boxStyle,
		Alignment: draw.TextAlignmentRight,
	})
}

func formatCompact(startTime time.Time) string {
	if startTime.IsZero() {
		return "N/A"
	}

	d := time.Since(startTime)

	if d.Hours() >= 24 {
		days := int(d.Hours() / 24)
		hours := int(d.Hours()) % 24
		return fmt.Sprintf("%dd %dh", days, hours)
	}

	return d.Round(time.Second).String()
}
