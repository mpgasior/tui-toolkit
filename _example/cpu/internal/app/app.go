package app

import (
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/model"
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/task"
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/ui"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/mvu"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
	"github.com/mpgasior/tui-toolkit/vt"
)

type App struct {
	state model.State
	ui    ui.ViewState
	store *process.Store
}

func New() *App {
	return &App{
		state: model.State{
			SortBy:    model.SortByRecentCPU,
			SortOrder: model.SortOrderDescending,
		},
		store: process.NewStore(),
		ui:    ui.New(),
	}
}

func (a *App) Init() mvu.Task {
	return task.Refresh(a.store, time.Second)
}

func (a *App) Update(e mvu.Event) mvu.Task {
	switch e.(type) {
	case task.DataRefreshedEvent:
		return a.TaskQuery()
	case task.QueryResultEvent:
		result := e.(task.QueryResultEvent)
		if synced := a.state.Sync(result.Data); synced {
			a.ui.Table.SortBy = a.state.SortBy
			a.ui.Table.SortOrder = a.state.SortOrder
			a.ui.Table.Rows = a.state.Rows
		}
		return a.TaskStopQuery()
	case task.TickEvent:
		a.ui.Search.Spinner.Next()
		return mvu.TaskNone
	}

	if key, ok := e.(vt.KeyEvent); ok {
		if key.IsKey(vt.KeyCtrlC) {
			return mvu.TaskShutdown
		}

		switch key.Key {
		case vt.KeyTab:
			a.ui.NextFocus()
			return mvu.TaskNone
		case vt.KeyShiftTab:
			a.ui.PrevFocus()
			return mvu.TaskNone
		}

		switch a.ui.CurrentFocus {
		case ui.FocusSearch:
			if key.IsKey(vt.KeyEnter) {
				a.ui.CurrentFocus = ui.FocusTable
				return mvu.TaskNone
			}
			if consumed := a.ui.Search.Update(key); consumed {
				a.state.SearchTerm = a.ui.Search.String()
				return a.TaskQuery()
			}
		case ui.FocusTable:
			switch key.Rune {
			case 's':
				a.state.SortOrder = (a.state.SortOrder + 1) % 2
				return a.TaskQuery()
			case 'h':
				a.state.SortBy = ui.PrevSortBy(a.state.SortBy)
				return a.TaskQuery()
			case 'l':
				a.state.SortBy = ui.NextSortBy(a.state.SortBy)
				return a.TaskQuery()
			}
		}
	}

	return mvu.TaskNone
}

func (a *App) TaskQuery() mvu.Task {
	a.ui.SetSearching(true)
	return mvu.TaskN(
		task.Tick(a.ui.Search.Spinner.ID, 80*time.Millisecond),
		task.Query(a.store, a.state.CurrentQuery()),
	)
}

func (a *App) TaskStopQuery() mvu.Task {
	a.ui.SetSearching(false)
	return mvu.TaskCancel(a.ui.Search.Spinner.ID)
}

func (a *App) Render(ctx mvu.RenderContext) {
	ctx.View.SetCursorPos(-1, -1)
	draw.Clear(ctx.View, screen.DefaultStyle)

	layout := view.SplitH(ctx.View,
		view.Fixed("search", 3),
		view.Dynamic("body", 3),
		view.Fixed("help", 1))
	search, table, help := layout["search"], layout["body"], layout["help"]

	a.ui.Search.Draw(search, a.ui.CurrentFocus == ui.FocusSearch)
	a.ui.Table.Draw(table, a.ui.CurrentFocus == ui.FocusTable)

	draw.Line(help, "Quit: ctrl+c | Focus: [Shift]Tab | Sort Order: s | Sort By: [h][l]", screen.DefaultStyle.Fg(screen.ColorBlue))
}
