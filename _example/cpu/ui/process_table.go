package ui

import (
	"strconv"

	"github.com/mpgasior/tui-toolkit/_example/cpu/process"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/mvu"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
)

type ProcessTable struct {
	Rows []process.ProcessInfo
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
		Attr(screen.AttrUnderline)

	layout := view.SplitH(body, view.Fixed("th", 1), view.Dynamic("tb", 1))
	tHead, tBody := layout["th"], layout["tb"]

	drawLine(tHead,
		draw.TextChunk{"", headerStyle},
		draw.TextChunk{"PID", headerStyle},
		draw.TextChunk{"Name", headerStyle},
		draw.TextChunk{"Kernel", headerStyle},
		draw.TextChunk{"[User]", headerStyle.Fg(screen.ColorGreen)})

	_, h := tBody.Size()

	for idx, info := range t.Rows {
		if idx >= h {
			break
		}

		row := tBody.Offset(idx, 0, 0, 0)

		drawInfo(row, info)
	}
}

func drawLine(vp view.Port,
	selected draw.TextChunk,
	pid draw.TextChunk,
	name draw.TextChunk,
	kernel draw.TextChunk,
	user draw.TextChunk) {
	layout := view.SplitV(vp,
		view.Fixed("selected", 4),
		view.Fixed("pid", 7),
		view.Dynamic("name", 25),
		view.Dynamic("kernel", 5),
		view.Dynamic("user", 5))

	draw.Text(layout["selected"], selected)
	draw.Text(layout["pid"], pid)
	draw.Text(layout["name"], name)
	draw.Text(layout["kernel"], kernel)
	draw.Text(layout["user"], user)
}

func drawInfo(vp view.Port, info process.ProcessInfo) {
	drawLine(vp,
		draw.TextChunk{
			Text:  "[ ]",
			Style: screen.DefaultStyle,
		},
		draw.TextChunk{
			Text:  strconv.FormatInt(int64(info.PID), 10),
			Style: screen.DefaultStyle,
		},
		draw.TextChunk{
			Text:  info.Name,
			Style: screen.DefaultStyle,
		},
		draw.TextChunk{
			Text:  info.KernelTime.String(),
			Style: screen.DefaultStyle,
		},
		draw.TextChunk{
			Text:  info.UserTime.String(),
			Style: screen.DefaultStyle,
		})
}
