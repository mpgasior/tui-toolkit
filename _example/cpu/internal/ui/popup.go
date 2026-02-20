package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
)

type Popup struct {
	PID    uint32
	Loaded bool
	//Result model.QueryResult
	Result process.Profile
}

func (p *Popup) Reset() {
	p.Loaded = false
}

func (p *Popup) Draw(vp view.Port) {
	draw.Clear(vp, screen.DefaultStyle)
	draw.Box(vp, draw.BoxBorderDouble, screen.DefaultStyle.Fg(screen.ColorGreen))

	if !p.Loaded {
		pid := strconv.FormatInt(int64(p.PID), 10)
		draw.Line(vp.Offset(1), "Loading PID ("+pid+") ...", screen.DefaultStyle)
		return
	}

	mainLayout := view.SplitH(vp.Offset(1),
		view.Fixed("details", 7),
		view.Dynamic("chart", 1),
		view.Fixed("help", 1),
	)

	detailsForm := view.SplitH(mainLayout["details"].Offset(0, 0, 0, 1),
		view.Fixed("pid", 1),
		view.Fixed("name", 1),
		view.Fixed("creation-time", 1),
		view.Fixed("avg-cpu", 1),
		view.Fixed("recent-cpu", 1),
		view.Fixed("peak-mem", 1),
		view.Fixed("recent-mem", 1),
	)

	setField := func(key string, title string, value string) {
		fieldLayout := view.SplitV(detailsForm[key],
			view.Fixed("title", 14),
			view.Fixed("gap", 1),
			view.Dynamic("value", 1),
		)

		draw.Line(fieldLayout["title"], title, screen.DefaultStyle.Attr(screen.AttrBold))
		draw.Line(fieldLayout["value"], value, screen.DefaultStyle)
	}

	stats, _ := p.Result.History.Stats()

	setField("pid", "PID", strconv.FormatInt(int64(p.Result.Info.PID), 10))
	setField("name", "Name", p.Result.Info.Name)
	if !p.Result.CreationTime.IsZero() {
		setField("creation-time", "Creation Time", p.Result.CreationTime.String())
	}
	setField("avg-cpu", "CPU% (Avg 1m)", fmt.Sprintf("%.2f%%", stats.AvgCPU))
	setField("recent-cpu", "CPU% (Now)", fmt.Sprintf("%.2f%%", stats.RecentCPU))
	setField("peak-mem", "MEM (Peak)", formatWorkingSet(stats.PeakWorkingSet))
	setField("recent-mem", "MEM (Now)", formatWorkingSet(stats.WorkingSet))

	for idx := range p.Result.History.Len() {
		sample := p.Result.History.Get(idx)
		v := strconv.FormatInt(int64(sample.WorkingSet), 10)
		draw.Line(mainLayout["chart"].Offset(idx, 0, 0, 0), v, screen.DefaultStyle)
	}

	items := []string{"Kill: k", "Interrupt: i", "Cancel: Esc"}
	text := strings.Join(items, " • ")
	draw.Line(mainLayout["help"], text, screen.DefaultStyle.Fg(screen.ColorBlue))
}
