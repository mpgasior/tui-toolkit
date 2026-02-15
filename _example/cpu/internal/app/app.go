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
		store: process.NewStore(),
		ui:    ui.New(),
	}
}

func (a *App) Init() mvu.Task {
	return tasks.TaskRefresh(a.store, time.Second)
}

func (a *App) Update(e mvu.Event) mvu.Task {
	switch e.(type) {
	case tasks.DataRefreshedEvent:
	case tasks.QueryResultEvent:
	case tasks.TickEvent:
	}

	if key, ok := e.(vt.KeyEvent); ok {
		if key.IsKey(vt.KeyCtrlC) {
			return mvu.TaskShutdown
		}

		switch key.Key {
		case vt.KeyTab:
			a.ui.NextFocus()
		case vt.KeyShiftTab:
			a.ui.PrevFocus()
		}

		if a.ui.CurrentFocus == ui.FocusSearch {
			if consumed := a.ui.TextInput.Update(key); consumed {
				a.state.SearchTerm = a.ui.TextInput.String()
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
	search, _, help := layout["search"], layout["body"], layout["help"]

	a.renderSearch(search)

	draw.Line(help, "[ctrl+c] Quit", screen.DefaultStyle)
}

func (a *App) renderSearch(vp view.Port) {
	boxStyle := screen.DefaultStyle
	if a.ui.CurrentFocus == ui.FocusSearch {
		boxStyle = boxStyle.Fg(screen.ColorGreen)
	}
	draw.Box(vp, draw.BoxBorderRounded, boxStyle)

	layout := view.SplitV(vp.Offset(1),
		view.Dynamic("input", 1),
		view.Fixed("spinner", 1),
	)

	inputView, spinnerView := layout["input"], layout["spinner"]

	w, _ := inputView.Size()
	text, cursor := a.ui.TextInput.Slice(w)
	if len(text) == 0 {
		draw.Line(inputView, "Search...", screen.DefaultStyle.Fg(screen.ColorBlue))
	} else {
		draw.Line(inputView, string(text), screen.DefaultStyle)
	}

	if a.ui.CurrentFocus == ui.FocusSearch {
		inputView.SetCursorPos(cursor, 0)
	}

	if a.state.IsLoading {
		r := a.ui.Spinner.Frame()
		draw.Rune(spinnerView, 0, 0, r, screen.DefaultStyle)
	}
}
