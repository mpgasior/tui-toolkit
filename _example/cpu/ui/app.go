package ui

import (
	"context"
	"time"

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
	store          *process.ProcessStore
	spinnerFrame   int
	querying       bool
}

func New() *App {
	return &App{
		input: &SearchInput{},
		table: &ProcessTable{},
		store: &process.ProcessStore{},
	}
}

type TickEvent struct{}

func TaskTick() mvu.Task {
	return mvu.Task{
		ID: "tick",
		Execute: func(ctx context.Context, ch chan<- mvu.Event) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(80 * time.Millisecond):
				}

				select {
				case <-ctx.Done():
					return
				case ch <- TickEvent{}:
				}
			}
		},
	}
}

func TaskCancelTick() mvu.Task {
	return mvu.TaskCancel("tick")
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
			a.querying = true
			return mvu.TaskN(
				TaskTick(),
				worker.TaskQuery(a.store, string(a.input.Term)),
			)
		case FocusTable:
		}

		return mvu.TaskNone
	}

	if _, ok := e.(TickEvent); ok {
		a.spinnerFrame += 1
		return mvu.TaskNone
	}

	if r, ok := e.(worker.QueryResultEvent); ok {
		a.spinnerFrame = 0
		a.querying = false
		a.table.Rows = r.Rows
		return TaskCancelTick()
	}

	if _, ok := e.(worker.DataRefreshedEvent); ok {
		a.querying = true
		return mvu.TaskN(
			TaskTick(),
			worker.TaskQuery(a.store, string(a.input.Term)),
		)
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

	if a.querying {
		draw.Spinner(search.Offset(1, 0, 0, 10), a.spinnerFrame, draw.SpinnerBraille, screen.DefaultStyle)
	}

	draw.Line(help, "[ctrl+c] Quit", screen.DefaultStyle)
}
