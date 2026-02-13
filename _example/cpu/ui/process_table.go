package ui

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/mpgasior/tui-toolkit/_example/cpu/process"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/mvu"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
)

type ProcessTable struct {
	Rows []process.Profile
}

func (t *ProcessTable) Render(ctx mvu.RenderContext) {
	draw.Clear(ctx.View, screen.DefaultStyle)

	boxStyle := screen.DefaultStyle
	if ctx.Focused {
		boxStyle = boxStyle.Fg(screen.ColorGreen)
	}
	draw.Box(ctx.View, draw.BoxBorderThin, boxStyle)

	body := ctx.View.Offset(1)

	if len(t.Rows) == 0 {
		center := view.Center(body, view.Dynamic("w", 1), view.Dynamic("h", 1))
		draw.Line(center, "waiting...", screen.DefaultStyle)
		return
	}

	headerStyle := screen.DefaultStyle.
		Attr(screen.AttrUnderline | screen.AttrBold)

	layout := view.SplitH(body, view.Fixed("th", 1), view.Dynamic("tb", 1))
	tHead, tBody := layout["th"], layout["tb"]

	drawLine(tHead,
		draw.TextChunk{Text: "", Style: headerStyle},
		draw.TextChunk{Text: "PID", Style: headerStyle},
		draw.TextChunk{Text: "Name", Style: headerStyle},
		draw.TextChunk{Text: "CPU% (Avg)", Style: headerStyle, Alignment: draw.TextAlignmentRight},
		draw.TextChunk{Text: "CPU% (Now)", Style: headerStyle.Fg(screen.ColorGreen)})

	_, h := tBody.Size()

	for idx, info := range t.Rows {
		if idx >= h {
			break
		}

		row := tBody.Offset(idx, 0, 0, 0)

		drawInfo(row, info)
	}

	w, h := ctx.View.Size()
	text := fmt.Sprintf("%d of %d", 0, len(t.Rows))

	pos := w - 1 - utf8.RuneCountInString(text)

	draw.Line(ctx.View.Slice(pos, h-1, w, h), text, screen.DefaultStyle)
}

func drawLine(vp view.Port,
	selected draw.TextChunk,
	pid draw.TextChunk,
	name draw.TextChunk,
	cpuAvg draw.TextChunk,
	cpuRecent draw.TextChunk) {
	layout := view.SplitV(vp,
		view.Fixed("selected", 4),
		view.Fixed("pid", 7),
		view.Fixed("avg-cpu", 10),
		view.Fixed("", 2),
		view.Fixed("recent-cpu", 10),
		view.Fixed("", 2),
		view.Dynamic("name", 15))

	draw.Text(layout["selected"], selected)
	draw.Text(layout["pid"], pid)
	draw.Text(layout["name"], name)
	draw.Text(layout["avg-cpu"], cpuAvg)
	draw.Text(layout["recent-cpu"], cpuRecent)
}

func drawInfo(vp view.Port, p process.Profile) {
	stats := p.History.Stats()

	drawLine(vp,
		draw.TextChunk{
			Text:  "[ ]",
			Style: screen.DefaultStyle,
		},
		draw.TextChunk{
			Text:  strconv.FormatInt(int64(p.Info.PID), 10),
			Style: screen.DefaultStyle,
		},
		draw.TextChunk{
			Text:  p.Info.Name,
			Style: screen.DefaultStyle,
		},
		draw.TextChunk{
			Text: func() string {
				if p.History == nil || p.History.Len() < 2 {
					return ""
				}
				return fmt.Sprintf("%5.2f%%", stats.AvgCPU*100)
			}(),
			Style:     screen.DefaultStyle,
			Alignment: draw.TextAlignmentRight,
		},
		draw.TextChunk{
			Text: func() string {
				if p.History == nil || p.History.Len() < 2 {
					return ""
				}
				return fmt.Sprintf("%5.2f%%", stats.RecentCPU*100)
			}(),
			Style:     screen.DefaultStyle,
			Alignment: draw.TextAlignmentRight,
		})
}
