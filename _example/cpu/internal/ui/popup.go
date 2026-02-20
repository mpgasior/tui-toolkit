package ui

import (
	"strconv"
	"strings"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/model"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
)

type Popup struct {
	Result model.QueryResult
}

func (p *Popup) Draw(vp view.Port) {
	draw.Clear(vp, screen.DefaultStyle)
	draw.Box(vp, draw.BoxBorderDouble, screen.DefaultStyle.Fg(screen.ColorGreen))

	mainLayout := view.SplitH(vp.Offset(1),
		view.Fixed("details", 7),
		view.Dynamic("dynamic", 1),
		view.Fixed("help", 1),
	)

	detailsForm := view.SplitH(mainLayout["details"],
		view.Fixed("pid", 1),
		view.Fixed("name", 1),
		view.Fixed("age", 1),
		view.Fixed("avg-cpu", 1),
		view.Fixed("recent-cpu", 1),
		view.Fixed("peak-mem", 1),
		view.Fixed("recent-mem", 1),
	)

	setField := func(key string, title string, value string) {
		fieldLayout := view.SplitV(detailsForm[key],
			view.Fixed("title", 15),
			view.Dynamic("value", 1),
		)

		draw.Line(fieldLayout["title"], title, screen.DefaultStyle.Attr(screen.AttrBold))
		draw.Line(fieldLayout["value"], value, screen.DefaultStyle)
	}

	setField("pid", "PID: ", strconv.FormatInt(int64(p.Result.PID), 10))
	setField("name", "Name: ", p.Result.Name)
	setField("age", "Age: ", p.Result.Age.String())
	setField("avg-cpu", "CPU% (Avg 1m): ", formatPercentage(p.Result.AvgCPU))
	setField("recent-cpu", "CPU% (Now): ", formatPercentage(p.Result.RecentCPU))
	setField("peak-mem", "MEM (Peak): ", formatWorkingSet(p.Result.PeakWorkingSet))
	setField("recent-mem", "MEM (Now): ", formatWorkingSet(p.Result.WorkingSet))

	draw.Line(vp.Offset(0, 0, 0, 5), p.Result.Name, screen.DefaultStyle)

	items := []string{"Kill: k", "Interrupt: i", "Cancel: Esc"}
	text := strings.Join(items, " • ")
	draw.Line(mainLayout["help"], text, screen.DefaultStyle.Fg(screen.ColorBlue))
}
