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
}

func New() *App {
	return &App{
		state: model.State{
			Store:     process.NewStore(),
			SortBy:    model.SortByRecentCPU,
			SortOrder: model.SortOrderDescending,
		},
		ui: ui.New(),
	}
}

func (a *App) Init() mvu.Task {
	return task.Refresh(a.state.Store, time.Second)
}

func (a *App) Update(e mvu.Event) mvu.Task {
	switch event := e.(type) {
	case task.DataRefreshedEvent:
		if a.ui.IsFocused(ui.FocusPopup) {
			return task.QuerySingle(a.state.Store, a.state.PID)
		}
		return a.TaskQuery()
	case task.QuerySingleResultEvent:
		a.ui.Popup.Result = event.Result
		a.ui.Popup.Loaded = true
		return mvu.TaskNone
	case task.QueryResultEvent:
		a.state.Sync(event.Data)

		a.ui.Table.SortBy = a.state.SortBy
		a.ui.Table.SortOrder = a.state.SortOrder
		a.ui.Table.Rows = a.state.Filtered

		return a.TaskStopQuery()
	case task.TickEvent:
		a.ui.Search.Spinner.Next()
		return mvu.TaskNone
	}

	if _, ok := e.(vt.PasteEvent); ok && a.ui.IsFocused(ui.FocusSearch) {
		if didUpdate := a.ui.Search.Update(e); didUpdate {
			a.state.SearchTerm = a.ui.Search.String()
			return a.TaskQuery()
		}
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
		case vt.KeyCtrlP:
			isPaused := a.state.TogglePause()
			a.ui.Table.IsPaused = isPaused
			if isPaused {
				return task.CancelRefresh()
			}
			return task.Refresh(a.state.Store, time.Second)
		}

		switch a.ui.CurrentFocus {
		case ui.FocusSearch:
			if key.IsKey(vt.KeyEnter) {
				a.ui.CurrentFocus = ui.FocusTable
				return mvu.TaskNone
			}
			if didUpdate := a.ui.Search.Update(e); didUpdate {
				a.state.SearchTerm = a.ui.Search.String()
				return a.TaskQuery()
			}
		case ui.FocusTable:
			if key.IsKey(vt.KeyEsc) {
				a.ui.Table.Reset()
				a.ui.CurrentFocus = ui.FocusSearch
			}
			if didUpdate := a.ui.Table.Update(key); didUpdate {
				a.state.SortBy = a.ui.Table.SortBy
				a.state.SortOrder = a.ui.Table.SortOrder
				if pid, ok := a.ui.Table.GetPID(); ok {
					a.state.PID = pid
					a.ui.Popup.PID = a.state.PID

					a.ui.CurrentFocus = ui.FocusPopup
					return task.QuerySingle(a.state.Store, a.state.PID)
				}

				return a.TaskQuery()
			}
		case ui.FocusPopup:
			if key.IsKey(vt.KeyEsc) {
				a.ui.CurrentFocus = ui.FocusTable
				a.ui.Popup.Reset()
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
		task.Query(a.state.Store, a.state.CurrentQuery()),
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
		view.Dynamic("table", 1),
		view.Fixed("help", 1),
	)

	a.ui.Search.Draw(layout["search"], a.ui.CurrentFocus == ui.FocusSearch)
	a.ui.Table.Draw(layout["table"], a.ui.CurrentFocus == ui.FocusTable)
	if a.ui.IsFocused(ui.FocusPopup) {
		a.ui.Popup.Draw(ctx.View.Offset(5, 10, 5, 10))
	}

	ui.DrawHelp(layout["help"], a.ui.CurrentFocus)
}
