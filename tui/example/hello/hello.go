package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mpgasior/tui-go/tui"
	"github.com/mpgasior/tui-go/vt"
)

type Ticker struct {
	maxTicks int
	ticks    int
}

type TickEvent struct{}

func (t *Ticker) Init() tui.Task {
	f := func(ctx context.Context, ch chan<- tui.Event) {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				e := TickEvent{}
				select {
				case <-ctx.Done():
					return
				case ch <- e:
				}
			}
		}
	}

	return tui.TaskF(f)
}

func (t *Ticker) Update(e tui.Event) tui.Task {
	if _, ok := e.(TickEvent); ok {
		t.ticks += 1
	}

	if t.ticks == t.maxTicks {
		return tui.TaskShutdown
	}

	if e, ok := e.(tui.KeyEvent); ok {
		if e.IsKey(vt.KeyCtrlC, vt.KeyEsc) {
			return tui.TaskShutdown
		}

		if e.IsRune('q') {
			return tui.TaskShutdown
		}
	}
	return tui.TaskNone
}

func (t *Ticker) Render(ctx tui.RenderContext) {
	left := t.maxTicks - t.ticks
	text := fmt.Sprintf("Stopping in %ds...", left)
	style := tui.DefaultStyle.
		Fg(tui.ColorHex(0xFFA500)).
		Bg(tui.ColorHex(0x303030))

	for i, r := range []rune(text) {
		ctx.Viewport.SetAt(i, 0, r, nil, 1, style)
	}
}

func main() {
	ticker := &Ticker{
		maxTicks: 10,
	}

	if err := tui.Run(ticker); err != nil {
		log.Fatalf("%v", err)
	}
}
