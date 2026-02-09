package ui

import (
	"github.com/mpgasior/tui-toolkit/_example/cpu/process"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
	"github.com/mpgasior/tui-toolkit/vt"
)

type UI struct {
	input *SearchInput
	table *ProcessTable
}

func New() *UI {
	return &UI{
		input: &SearchInput{},
		table: &ProcessTable{},
	}
}

func (a *UI) Update(e vt.KeyEvent) {
	a.input.Update(e)
}

func (a *UI) OnNewRows(r []process.ProcessInfo) {
	a.table.Rows = r
}

func (a *UI) Draw(m screen.Mutator) {
	vp := view.NewPort(m)

	draw.Clear(vp, screen.DefaultStyle)

	layout := view.SplitH(vp, view.Fixed("search", 3), view.Dynamic("body", 3), view.Fixed("help", 1))
	search, body, help := layout["search"], layout["body"], layout["help"]

	a.input.Draw(search)
	a.table.Draw(body)

	draw.Line(help, "[Q]uit", screen.DefaultStyle)
}
