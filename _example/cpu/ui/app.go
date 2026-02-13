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
	spinner        *Spinner
	store          *process.Store
	querying       bool
	rows           []worker.QueryResult
	sortByIndex    int
}

var columns = []string{"RecentCPU", "PID", "AvgCPU", "Name"}

func New() *App {
	return &App{
		input: &SearchInput{},
		spinner: &Spinner{
			ID:    "spinner",
			Frame: draw.SpinnerBraille,
		},
		store: process.NewStore(),
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

		if key.IsKey(vt.KeyF4) {
			a.sortByIndex = (a.sortByIndex + 1) % len(columns)
			return a.QueryTask()
		}

		if key.IsKey(vt.KeyTab) {
			a.focusedElement = (a.focusedElement + 1) % 2
			return mvu.TaskNone
		} else if key.IsKey(vt.KeyShiftTab) {
			a.focusedElement = (a.focusedElement - 1 + 2) % 2
			return mvu.TaskNone
		}

		switch a.focusedElement {
		case FocusSearch:
			a.input.Update(e)
			return a.QueryTask()
		case FocusTable:
		}

		return mvu.TaskNone
	}

	if t, ok := e.(SpinnerTickEvent); ok {
		a.spinner.OnTick(t)
		return mvu.TaskNone
	}

	if r, ok := e.(worker.QueryResultEvent); ok {
		a.rows = r.Rows
		return a.StopSpinnerTask()
	}

	if _, ok := e.(worker.DataRefreshedEvent); ok {
		return a.QueryTask()
	}

	return mvu.TaskNone
}

func (a *App) QueryTask() mvu.Task {
	a.querying = true

	query := worker.Query{
		Term:   string(a.input.Term),
		SortBy: columns[a.sortByIndex],
	}
	return mvu.TaskN(
		a.spinner.StartTask(),
		worker.TaskQuery(a.store, query),
	)
}

func (a *App) StopSpinnerTask() mvu.Task {
	a.querying = false
	return a.spinner.CancelTask()
}

func (a *App) Render(ctx mvu.RenderContext) {
	ctx.View.SetCursorPos(-1, -1)
	draw.Clear(ctx.View, screen.DefaultStyle)

	layout := view.SplitH(ctx.View,
		view.Fixed("search", 3),
		view.Dynamic("body", 3),
		view.Fixed("help", 1))
	search, body, help := layout["search"], layout["body"], layout["help"]

	a.input.Render(mvu.RenderContext{
		View:    search,
		Focused: a.focusedElement == FocusSearch,
	})

	RenderProcessTable(mvu.RenderContext{
		View:    body,
		Focused: a.focusedElement == FocusTable},
		a.rows,
		columns[a.sortByIndex])

	if a.querying {
		w, _ := search.Size()
		a.spinner.Render(search.Offset(1, 0, 0, w-2))
	}

	draw.Line(help, "[ctrl+c] Quit", screen.DefaultStyle)
}
