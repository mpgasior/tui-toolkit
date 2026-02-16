package app

import (
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/model"
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/tasks"
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
	a.ui.SetSearching(true)
	return tasks.TaskRefresh(a.store, time.Second)
}

func (a *App) Update(e mvu.Event) mvu.Task {
	switch e.(type) {
	case tasks.DataRefreshedEvent:
		a.ui.SetSearching(true)
		return mvu.TaskN(
			tasks.TaskTick(a.ui.Search.Spinner.ID, 80*time.Millisecond),
			tasks.TaskQuery(a.store, a.state.CurrentQuery()),
		)
	case tasks.QueryResultEvent:
		a.ui.SetSearching(false)
		data := e.(tasks.QueryResultEvent).Data
		if synced := a.state.Sync(data); synced {
			a.ui.Table.Rows = data
		}
		return mvu.TaskCancel(a.ui.Search.Spinner.ID)
	case tasks.TickEvent:
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

		if a.ui.IsFocused(ui.FocusSearch) {
			if consumed := a.ui.Search.Update(key); consumed {
				a.state.SearchTerm = a.ui.Search.String()
			}
		}
	}

	return mvu.TaskNone
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

	draw.Line(help, "[ctrl+c] Quit", screen.DefaultStyle)
}
