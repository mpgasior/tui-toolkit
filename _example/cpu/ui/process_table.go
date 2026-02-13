package ui

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/mpgasior/tui-toolkit/_example/cpu/worker"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/mvu"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
)

type ProcessTable struct {
	Rows []worker.QueryResult
}

func RenderProcessTable(
	ctx mvu.RenderContext,
	rows []worker.QueryResult,
	sortBy string) {
	draw.Clear(ctx.View, screen.DefaultStyle)

	boxStyle := screen.DefaultStyle
	if ctx.Focused {
		boxStyle = boxStyle.Fg(screen.ColorGreen)
	}
	draw.Box(ctx.View, draw.BoxBorderThin, boxStyle)

	body := ctx.View.Offset(1)

	if len(rows) == 0 {
		center := view.Center(body, view.Dynamic("w", 1), view.Dynamic("h", 1))
		draw.Line(center, "waiting...", screen.DefaultStyle)
		return
	}

	headerStyle := screen.DefaultStyle.
		Attr(screen.AttrUnderline | screen.AttrBold)

	layout := view.SplitH(body, view.Fixed("th", 1), view.Dynamic("tb", 1))
	tHead, tBody := layout["th"], layout["tb"]

	m := map[string]*draw.TextChunk{
		"PID":       &draw.TextChunk{Text: "PID", Style: headerStyle},
		"Name":      &draw.TextChunk{Text: "Name", Style: headerStyle},
		"AvgCPU":    &draw.TextChunk{Text: "CPU% (Avg 1m)", Style: headerStyle},
		"RecentCPU": &draw.TextChunk{Text: "CPU% (Now)", Style: headerStyle},
		"Memory":    &draw.TextChunk{Text: "Memory", Style: headerStyle},
	}

	style := m[sortBy].Style
	m[sortBy].Style = style.Fg(screen.ColorGreen)

	drawLine(tHead,
		draw.TextChunk{Text: "", Style: headerStyle},
		*m["PID"],
		*m["Name"],
		*m["AvgCPU"],
		*m["RecentCPU"],
		*m["Memory"])
	_, h := tBody.Size()

	for idx, info := range rows {
		if idx >= h {
			break
		}

		row := tBody.Offset(idx, 0, 0, 0)

		drawInfo(row, info)
	}

	w, h := ctx.View.Size()
	text := fmt.Sprintf("%d of %d", 0, len(rows))

	pos := w - 1 - utf8.RuneCountInString(text)

	draw.Line(ctx.View.Slice(pos, h-1, w, h), text, screen.DefaultStyle)
}

func drawLine(vp view.Port,
	selected draw.TextChunk,
	pid draw.TextChunk,
	name draw.TextChunk,
	cpuAvg draw.TextChunk,
	cpuRecent draw.TextChunk,
	memory draw.TextChunk) {
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

	draw.Text(layout["selected"], selected)
	draw.Text(layout["pid"], pid)
	draw.Text(layout["avg-cpu"], cpuAvg)
	draw.Text(layout["recent-cpu"], cpuRecent)
	draw.Text(layout["memory"], memory)
	draw.Text(layout["name"], name)
}

func drawInfo(vp view.Port, r worker.QueryResult) {
	drawLine(vp,
		draw.TextChunk{
			Text:  "[ ]",
			Style: screen.DefaultStyle,
		},
		draw.TextChunk{
			Text:  strconv.FormatInt(int64(r.PID), 10),
			Style: screen.DefaultStyle,
		},
		draw.TextChunk{
			Text:  r.Name,
			Style: screen.DefaultStyle,
		},
		draw.TextChunk{
			Text:      fmt.Sprintf("%5.2f%%", r.AvgCPU*100),
			Style:     screen.DefaultStyle,
			Alignment: draw.TextAlignmentRight,
		},
		draw.TextChunk{
			Text:      fmt.Sprintf("%5.2f%%", r.RecentCPU*100),
			Style:     screen.DefaultStyle,
			Alignment: draw.TextAlignmentRight,
		},
		draw.TextChunk{
			Text:      "?? MB",
			Style:     screen.DefaultStyle,
			Alignment: draw.TextAlignmentRight,
		})
}
