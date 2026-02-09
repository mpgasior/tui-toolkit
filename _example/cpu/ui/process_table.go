package ui

import (
	"strconv"

	"github.com/mpgasior/tui-toolkit/_example/cpu/process"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
)

type ProcessTable struct {
	Rows []process.ProcessInfo
}

func (pl *ProcessTable) Draw(v view.Port, focused bool) {
	draw.Clear(v, screen.DefaultStyle)
	boxStyle := screen.DefaultStyle
	if focused {
		boxStyle = boxStyle.Fg(screen.ColorGreen)
	}

	draw.Box(v, draw.BoxBorderDouble, boxStyle)

	v = v.Offset(1)

	if len(pl.Rows) == 0 {
		body := view.Center(v, view.Dynamic("w", 1), view.Dynamic("h", 1))
		draw.Line(body, "waiting...", screen.DefaultStyle)
		return
	}

	headerStyle := screen.DefaultStyle.
		Attr(screen.AttrUnderline)

	drawLine(v.Offset(1, 0, 0, 1),
		draw.TextChunk{"", headerStyle},
		draw.TextChunk{"PID", headerStyle},
		draw.TextChunk{"Name", headerStyle},
		draw.TextChunk{"Kernel", headerStyle},
		draw.TextChunk{"[User]", headerStyle.Fg(screen.ColorGreen)})

	tBody := v.Offset(2, 0, 0, 1)
	w, h := v.Size()

	for idx, info := range pl.Rows {
		if idx >= h-1 {
			break
		}

		row := tBody.Offset(idx, 0, 0, 0).Slice(0, 0, w, 1)

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
