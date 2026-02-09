package ui

import (
	"github.com/mpgasior/tui-toolkit/_example/cpu/process"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
	"github.com/mpgasior/tui-toolkit/vt"
)

type Focus int

const (
	FocusSearch Focus = iota
	FocusTable
)

type App struct {
	focusedElement Focus
	input          *SearchInput
	table          *ProcessTable
}

func New() *App {
	return &App{
		input: &SearchInput{},
		table: &ProcessTable{},
	}
}

func (a *App) HandleKey(e vt.KeyEvent) {
	if e.IsKey(vt.KeyTab) {
		a.focusedElement = FocusTable
	} else if e.IsKey(vt.KeyShiftTab) {
		a.focusedElement = FocusSearch
	}

	switch a.focusedElement {
	case FocusSearch:
		a.input.Update(e)
	case FocusTable:
	}
}

func (a *App) OnNewRows(r []process.ProcessInfo) {
	a.table.Rows = r
}

func (a *App) Draw(m screen.Mutator) {
	vp := view.NewPort(m)

	draw.Clear(vp, screen.DefaultStyle)

	layout := view.SplitH(vp, view.Fixed("search", 3), view.Dynamic("body", 3), view.Fixed("help", 1))
	search, body, help := layout["search"], layout["body"], layout["help"]

	a.input.Draw(search, a.focusedElement == FocusSearch)
	a.table.Draw(body, a.focusedElement == FocusTable)

	draw.Line(help, "[Q]uit", screen.DefaultStyle)
}
