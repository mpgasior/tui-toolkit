package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mpgasior/tui-go/tui"
	"github.com/mpgasior/tui-go/vt"
)

type Hello struct {
	timeout time.Duration
	left    time.Duration
}

type HelloTickEvent struct {
	left time.Duration
}

func (h *Hello) Init() tui.Task {
	f := func(ctx context.Context, ch chan<- tui.Event) {
		shutdownTime := time.Now().Add(h.timeout)
		shutdown := time.After(h.timeout)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-shutdown:
				select {
				case ch <- tui.ShutdownEvent:
				case <-ctx.Done():
					return
				}
			case <-ticker.C:
				e := HelloTickEvent{
					left: time.Until(shutdownTime),
				}
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

func (h *Hello) Update(e tui.Event) tui.Task {
	if e, ok := e.(HelloTickEvent); ok {
		h.left = e.left
	}

	if e, ok := e.(tui.KeyEvent); ok {
		if e.IsKey(vt.KeyCtrlC, vt.KeyEsc) {
			return tui.TaskShutdown
		}

		if e.IsRune('q') {
			return tui.TaskShutdown
		}
	}
	return tui.TaskNone()
}

func (h *Hello) Render(ctx tui.RenderContext) {
	left := h.left
	if h.left == time.Duration(0) {
		left = h.timeout
	}

	if left < 0 {
		left = time.Duration(0)
	}

	text := fmt.Sprintf("Stopping in %s...", left)
	style := tui.DefaultStyle.
		Fg(tui.ColorHex(0xFFA500)).
		Bg(tui.ColorHex(0x303030))

	for i, r := range []rune(text) {
		ctx.Viewport.SetAt(i, 0, r, nil, 1, style)
	}
}

func main() {
	app := tui.App{}

	c := &Hello{
		timeout: 10 * time.Second,
	}

	if err := app.Run(c); err != nil {
		log.Fatalf("%v", err)
	}
}
