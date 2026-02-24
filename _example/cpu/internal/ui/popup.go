package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/model"
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
)

type Popup struct {
	Key    process.Key
	Loaded bool
	data   model.ProcessHistory
}

func (p *Popup) Open(key process.Key) {
	p.Loaded = false
	p.Key = key
}

func (p *Popup) Update(data model.ProcessHistory) {
	p.Loaded = true
	p.data = data
}

func (p *Popup) Close() {
	p.Key = process.KeyNone
	p.Loaded = false
}

func (p *Popup) Draw(vp view.Port) {
	draw.Clear(vp, screen.DefaultStyle)
	draw.Box(vp, draw.BoxBorderDouble, screen.DefaultStyle.Fg(screen.ColorGreen))

	if !p.Loaded {
		pid := p.Key.String()
		draw.Line(vp.Offset(1), "Loading PID ("+pid+") ...", screen.DefaultStyle)
		return
	}

	mainLayout := view.SplitH(vp.Offset(1),
		view.Fixed("details", 8),
		view.Dynamic("chart", 1),
		view.Fixed("help", 1),
	)

	detailsForm := view.SplitH(mainLayout["details"].Offset(0, 0, 0, 1),
		view.Fixed("pid", 1),
		view.Fixed("name", 1),
		view.Fixed("creation-time", 1),
		view.Fixed("exit-time", 1),
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

	setField("pid", "PID", strconv.FormatInt(int64(p.data.PID), 10))
	setField("name", "Name", p.data.Name)
	if !p.data.CreationTime.IsZero() {
		setField("creation-time", "Creation Time", p.data.CreationTime.String())
	}
	if !p.data.ExitTime.IsZero() {
		setField("exit-time", "Exit Time", p.data.ExitTime.String())
	} else {
		setField("exit-time", "Exit Time", "")
	}
	if p.data.CPUReady {
		setField("avg-cpu", "CPU% (Avg 1m)", fmt.Sprintf("%.2f%%", p.data.AvgCPU))
		setField("recent-cpu", "CPU% (Now)", fmt.Sprintf("%.2f%%", p.data.LatestCPU))
	}
	if p.data.MemReady {
		setField("peak-mem", "MEM (Peak)", formatWorkingSet(p.data.MaxMem))
		setField("recent-mem", "MEM (Now)", formatWorkingSet(p.data.LatestMem))
	}

	//for idx := range p.Data.History.Len() {
	//	sample := p.Data.History.Get(idx)
	//	v := strconv.FormatInt(int64(sample.WorkingSet), 10)
	//	draw.Line(mainLayout["chart"].Offset(idx, 0, 0, 0), v, screen.DefaultStyle)
	//}

	items := []string{"Kill: k", "Interrupt: i", "Cancel: Esc"}
	text := strings.Join(items, " • ")
	draw.Line(mainLayout["help"], text, screen.DefaultStyle.Fg(screen.ColorBlue))
}
