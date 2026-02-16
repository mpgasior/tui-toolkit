package ui

import (
	"fmt"
	"strconv"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/model"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
)

type Table struct {
	Rows []model.QueryResult
}

func (t *Table) Draw(vp view.Port, focused bool) {
	boxStyle := screen.DefaultStyle
	if focused {
		boxStyle = boxStyle.Fg(screen.ColorGreen)
	}
	draw.Box(vp, draw.BoxBorderThin, boxStyle)

	if len(t.Rows) == 0 {
		draw.Line(vp, "work in progress", screen.DefaultStyle)
		return
	}

	layout := view.SplitV(vp,
		view.Fixed("selected", 4),
		view.Fixed("pid", 7),
		view.Fixed("avg-cpu", 13),
		view.Fixed("", 2),
		view.Fixed("recent-cpu", 10),
		view.Fixed("", 2),
		view.Fixed("memory", 10),
		view.Fixed("", 2),
		view.Dynamic("name", 15))

	cell := func(n string, row int) view.Port {
		col := layout[n]
		w, _ := col.Size()
		return col.Slice(0, row, w, 1)
	}

	headerStyle := screen.DefaultStyle.
		Attr(screen.AttrBold | screen.AttrUnderline)

	draw.Text(cell("pid", 0), draw.TextChunk{
		Text:  "PID",
		Style: headerStyle,
	})
	draw.Text(cell("name", 0), draw.TextChunk{
		Text:  "Name",
		Style: headerStyle,
	})
	draw.Text(cell("avg-cpu", 0), draw.TextChunk{
		Text:  "CPU% (Avg 1m)",
		Style: headerStyle,
	})
	draw.Text(cell("recent-cpu", 0), draw.TextChunk{
		Text:  "CPU% (Now)",
		Style: headerStyle.Fg(screen.ColorGreen),
	})

	_, h := vp.Size()

	for idx, row := range t.Rows {
		if idx >= h {
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
		draw.Text(cell("avg-cpu", idx+1), draw.TextChunk{
			Text:      fmt.Sprintf("%5.2f%%", row.AvgCPU),
			Style:     screen.DefaultStyle,
			Alignment: draw.TextAlignmentRight,
		})
		draw.Text(cell("recent-cpu", idx+1), draw.TextChunk{
			Text:      fmt.Sprintf("%5.2f%%", row.AvgCPU),
			Style:     screen.DefaultStyle,
			Alignment: draw.TextAlignmentRight,
		})
	}
}
