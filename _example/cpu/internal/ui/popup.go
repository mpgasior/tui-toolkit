package ui

import (
	"fmt"
	"math"
	"slices"
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
	chartLayout := view.SplitH(
		layout["chart"],
		view.Dynamic("cpu", 1),
		view.Dynamic("mem", 1),
	)
	p.drawChart(chartLayout["cpu"], "CPU Usage (%)", p.data.CPU, screen.DefaultStyle.Fg(screen.ColorRed))
	memory := make([]float64, len(p.data.Mem))
	for idx, m := range p.data.Mem {
		memory[idx] = float64(m)
	}
	p.drawChart(chartLayout["mem"], "Memory Usage (B)", memory, screen.DefaultStyle.Fg(screen.ColorBlue))
	p.drawHelp(layout["help"])
}

func (p *Popup) drawChart(vp view.Port, title string, data []float64, style screen.Style) {
	draw.Box(vp, draw.BoxBorderCorners, style)
	draw.Box(vp, draw.BoxBorderThin, style)

	if len(data) == 0 {
		return
	}
	w, h := vp.Size()

	dataMin := slices.Min(data)
	dataMax := slices.Max(data)
	now := data[len(data)-1]
	draw.Text(vp.Offset(0, 0, 0, 2), draw.TextChunk{
		Text:      " " + title + " ",
		Style:     style,
		Alignment: draw.TextAlignmentLeft,
	})

	draw.Text(vp.Offset(0, 2, 0, 0), draw.TextChunk{
		Text:      fmt.Sprintf(" %.2f%% (now)", now),
		Style:     style,
		Alignment: draw.TextAlignmentRight,
	})

	draw.Text(vp.Slice(0, h-1, w-2, h), draw.TextChunk{
		Text:      fmt.Sprintf(" %.2f%% (min) %.2f%% (max)", dataMin, dataMax),
		Style:     style,
		Alignment: draw.TextAlignmentRight,
	})

	vp = vp.Offset(1)
	dataRange := dataMax - dataMin
	rangeMax := float64(h)

	var blocks = [9]rune{' ', '▂', '▃', '▄', '▅', '▆', '▇', '█', '█'}

	for x, d := range data {
		if x >= w {
			break
		}

		var v float64
		if dataRange == 0 {
			v = float64(h) * 0.5
		} else {
			v = ((d - dataMin) / dataRange) * rangeMax
		}

		fullBlocks := int(math.Floor(v))
		remainder := v - float64(fullBlocks)
		partialIdx := int(remainder * 8)

		for yOffset := 0; yOffset < fullBlocks && yOffset < h; yOffset++ {
			draw.Rune(vp, x, h-1-yOffset, '█', style)
		}

		if fullBlocks < h {
			char := blocks[partialIdx]
			if char != ' ' {
				draw.Rune(vp, x, h-1-fullBlocks, char, style)
			}
		}
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
		view.Fixed("avg-cpu", 1),
		view.Fixed("recent-cpu", 1),
		view.Fixed("peak-mem", 1),
		view.Fixed("recent-mem", 1),
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

	if p.data.CPUReady {
		setField("avg-cpu", "CPU% (Avg 1m)", fmt.Sprintf("%.2f%%", p.data.AvgCPU))
		setField("recent-cpu", "CPU% (Now)", fmt.Sprintf("%.2f%%", p.data.LatestCPU))
	} else {
		setField("avg-cpu", "CPU% (Avg 1m)", "")
		setField("recent-cpu", "CPU% (Now)", "")
	}
	if p.data.MemReady {
		setField("peak-mem", "MEM (Peak)", formatWorkingSet(p.data.MaxMem))
		setField("recent-mem", "MEM (Now)", formatWorkingSet(p.data.LatestMem))
	} else {
		setField("peak-mem", "MEM (Peak)", "")
		setField("recent-mem", "MEM (Now)", "")
	}
}
