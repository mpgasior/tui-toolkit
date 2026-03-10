package ui

import (
	"fmt"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/model"
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
)

type Popup struct {
	key    process.Key
	loaded bool
	data   model.ProcessHistory
}

func (p *Popup) Open(key process.Key) {
	p.loaded = false
	p.key = key
}

func (p *Popup) Update(data model.ProcessHistory) {
	p.loaded = true
	p.data = data
}

func (p *Popup) Close() {
	p.key = process.KeyNone
	p.loaded = false
}

func (p *Popup) Draw(vp view.Port) {
	draw.Clear(vp, screen.DefaultStyle)
	draw.Box(vp, draw.BoxBorderDouble, screen.DefaultStyle.Fg(screen.ColorGreen))

	if !p.loaded {
		pid := p.key.String()
		draw.Line(vp.Offset(1), "Loading PID ("+pid+") ...", screen.DefaultStyle)
		return
	}

	text := fmt.Sprintf(" [ %s (%d) ]", p.data.Name, p.data.PID)
	draw.Line(vp.Offset(0, 0, 0, 2), text, screen.DefaultStyle)

	layout := view.SplitH(
		vp.Offset(1, 2),
		view.Fixed("details", 2),
		view.Dynamic("chart", 1),
	)

	p.drawDetails(layout["details"])
	p.drawHistograms(layout["chart"])
}

func (p *Popup) drawHistograms(vp view.Port) {
	charts := view.SplitV(
		vp,
		view.Dynamic("cpu", 1),
		view.Dynamic("mem", 1),
	)

	cpuStyle := screen.DefaultStyle.Fg(screen.ColorRed)
	drawMetric(charts["cpu"], p.data.CPU,
		" CPU Usage (%) ",
		fmt.Sprintf("%.2f%% (now) ", p.data.LatestCPU),
		fmt.Sprintf("%.2f%% (min) %.2f%% (avg) %.2f%% (max) ", p.data.MinCPU, p.data.AvgCPU, p.data.MaxCPU),
		cpuStyle,
		func(v float64) float64 { return v },
	)

	memStyle := screen.DefaultStyle.Fg(screen.ColorBlue)
	drawMetric(charts["mem"], p.data.Mem,
		" Memory Usage (%) ",
		fmt.Sprintf("%s (now) ", formatWorkingSet(p.data.LatestMem)),
		fmt.Sprintf("%s (min) %s (max) ", formatWorkingSet(p.data.MinMem), formatWorkingSet(p.data.MaxMem)),
		memStyle,
		func(v uint64) float64 { return float64(v) },
	)
}

func drawMetric[T any](
	viewPort view.Port,
	data []T,
	title string,
	nowLabel string,
	statsLabel string,
	style screen.Style,
	transform func(T) float64,
) {
	draw.Box(viewPort, draw.BoxBorderThin, style)

	histVP := viewPort.Offset(1)

	displayData := data
	if w, _ := histVP.Size(); len(displayData) > w {
		displayData = displayData[len(displayData)-w:]
	}

	draw.Histogram(histVP, displayData, transform, style)

	drawChartBorder(
		viewPort.Offset(0, 2, 0, 2),
		title,
		nowLabel,
		statsLabel,
		style,
	)
}

func drawChartBorder(
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

func (p *Popup) drawDetails(vp view.Port) {
	form := view.SplitH(vp,
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

	var vpTitle view.Port
	var vpValue view.Port

	vpTitle, vpValue = field("creation-time")
	draw.Line(vpTitle, "Created", styleTitle)
	if !p.data.CreationTime.IsZero() {
		draw.Line(vpValue, p.data.CreationTime.String(), styleValue)
	}

	if !p.data.ExitTime.IsZero() {
		vpTitle, vpValue = field("exit-time")
		draw.Line(vpTitle, "Exited", styleTitle)
		draw.Line(vpValue, p.data.ExitTime.String(), styleValue.Fg(screen.ColorRed))
	}
}
