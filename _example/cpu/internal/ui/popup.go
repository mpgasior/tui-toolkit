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

	layout := view.SplitH(
		vp.Offset(1, 2),
		view.Fixed("details", 4),
		view.Dynamic("chart", 1),
		view.Fixed("help", 1),
	)

	p.drawDetails(layout["details"])
	charts := view.SplitV(
		layout["chart"],
		view.Dynamic("cpu", 1),
		view.Dynamic("mem", 1),
	)

	cpuStyle := screen.DefaultStyle.Fg(screen.ColorRed)
	draw.Box(charts["cpu"], draw.BoxBorderThin, cpuStyle)
	draw.Histogram(charts["cpu"].Offset(1), p.data.CPU, func(v float64) float64 { return v }, cpuStyle)
	p.drawChartBorder(
		charts["cpu"].Offset(0, 2, 0, 2),
		" CPU Usage (%) ",
		fmt.Sprintf("%.2f%% (now) ", p.data.LatestCPU),
		fmt.Sprintf("%.2f%% (min) %.2f%% (avg) %.2f%% (max) ", p.data.MinCPU, p.data.AvgCPU, p.data.MaxCPU),
		cpuStyle,
	)

	memStyle := screen.DefaultStyle.Fg(screen.ColorBlue)
	draw.Box(charts["mem"], draw.BoxBorderThin, memStyle)
	draw.Histogram(charts["mem"].Offset(1), p.data.Mem, func(v uint64) float64 { return float64(v) }, memStyle)
	p.drawChartBorder(
		charts["mem"].Offset(0, 2, 0, 2),
		" Memory Usage (%) ",
		fmt.Sprintf("%s (now) ", formatWorkingSet(p.data.LatestMem)),
		fmt.Sprintf("%s (min) %s (max) ", formatWorkingSet(p.data.MinMem), formatWorkingSet(p.data.MaxMem)),
		memStyle,
	)
	p.drawHelp(layout["help"])
}

func (p *Popup) drawChartBorder(
	vp view.Port,
	topLeft, topRight, bottomRight string,
	style screen.Style,
) {
	if topLeft != "" {
		draw.Text(vp, draw.TextChunk{
			Text:  topLeft,
			Style: style,
		})
	}
	if topRight != "" {
		draw.Text(vp, draw.TextChunk{
			Text:      topRight,
			Style:     style,
			Alignment: draw.TextAlignmentRight,
		})
	}
	if bottomRight != "" {
		w, h := vp.Size()
		draw.Text(vp.Slice(0, h-1, w, h), draw.TextChunk{
			Text:      bottomRight,
			Style:     style,
			Alignment: draw.TextAlignmentRight,
		})
	}
}

func (p *Popup) drawHelp(vp view.Port) {
	items := []string{"Kill: k", "Interrupt: i", "Cancel: Esc"}
	text := strings.Join(items, " • ")
	draw.Line(vp, text, screen.DefaultStyle.Fg(screen.ColorBlue))
}

func (p *Popup) drawDetails(vp view.Port) {
	form := view.SplitH(vp,
		view.Fixed("pid", 1),
		view.Fixed("name", 1),
		view.Fixed("creation-time", 1),
		view.Fixed("exit-time", 1),
	)

	field := func(key string) (title view.Port, value view.Port) {
		layout := view.SplitV(form[key],
			view.Fixed("title", 14),
			view.Fixed("gap", 1),
			view.Dynamic("value", 1),
		)

		return layout["title"], layout["value"]
	}

	styleTitle := screen.DefaultStyle.Attr(screen.AttrBold)
	styleValue := screen.DefaultStyle

	setField := func(key, title, value string) {
		vpTitle, vpValue := field(key)

		draw.Line(vpTitle, title, styleTitle)
		draw.Line(vpValue, value, styleValue)
	}

	setField("pid", "PID", strconv.FormatInt(int64(p.data.PID), 10))
	setField("name", "Name", p.data.Name)

	var vpTitle view.Port
	var vpValue view.Port

	vpTitle, vpValue = field("creation-time")
	draw.Line(vpTitle, "Created", styleTitle)
	if !p.data.CreationTime.IsZero() {
		draw.Line(vpValue, p.data.CreationTime.String(), styleValue)
	}

	vpTitle, vpValue = field("exit-time")
	draw.Line(vpTitle, "Exited", styleTitle)
	if !p.data.ExitTime.IsZero() {
		draw.Line(vpValue, p.data.ExitTime.String(), styleValue.Fg(screen.ColorRed))
	}
}
