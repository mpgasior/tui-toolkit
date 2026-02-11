package ui

import (
	"github.com/mpgasior/tui-toolkit/_example/cpu/process"
	"github.com/mpgasior/tui-toolkit/_example/cpu/worker"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/mvu"
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
	snapshotID     int
	store          *process.ProcessStore
}

func New() *App {
	return &App{
		input: &SearchInput{},
		table: &ProcessTable{},
		store: &process.ProcessStore{},
	}
}

func (a *App) Init() mvu.Task {
	return worker.TaskRefresh(a.store)
}

func (a *App) Update(e mvu.Event) mvu.Task {
	if key, ok := e.(vt.KeyEvent); ok {
		if key.IsKey(vt.KeyCtrlC) {
			return mvu.TaskShutdown
		}

		if key.IsKey(vt.KeyTab) {
			a.focusedElement = FocusTable
		} else if key.IsKey(vt.KeyShiftTab) {
			a.focusedElement = FocusSearch
		}

		switch a.focusedElement {
		case FocusSearch:
			a.input.Update(e)
			a.snapshotID += 1
			return worker.TaskQuery(a.store, string(a.input.Term))
		case FocusTable:
		}

		return mvu.TaskNone
	}

	if r, ok := e.(worker.QueryResultEvent); ok {
		a.table.Rows = r.Rows
	}

	if _, ok := e.(worker.DataRefreshedEvent); ok {
		return worker.TaskQuery(a.store, string(a.input.Term))
	}

	return mvu.TaskNone
}

func (a *App) Render(ctx mvu.RenderContext) {
	vp := ctx.View
	draw.Clear(vp, screen.DefaultStyle)

	layout := view.SplitH(vp, view.Fixed("search", 3), view.Dynamic("body", 3), view.Fixed("help", 1))
	search, body, help := layout["search"], layout["body"], layout["help"]

	a.input.Render(mvu.RenderContext{
		View:    search,
		Focused: a.focusedElement == FocusSearch,
	})
	a.table.Render(mvu.RenderContext{
		View:    body,
		Focused: a.focusedElement == FocusTable,
	})

	draw.Line(help, "[ctrl+c] Quit", screen.DefaultStyle)
}
