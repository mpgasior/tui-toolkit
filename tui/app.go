package tui

import (
	"context"

	"github.com/mpgasior/tui-go/termx"
)

type App struct {
}

func (a App) Run(c Component) error {
	tty, err := termx.OpenTTY()
	if err != nil {
		return err
	}
	defer tty.Close()

	terminal, err := termx.New(tty.In, tty.Out)
	if err != nil {
		return err
	}

	w, h, err := terminal.GetSize()
	if err != nil {
		return err
	}

	ch := make(chan Event)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runtime := NewRuntime(ch)
	dispatch, shutdown := runtime.Start(ctx)
	defer shutdown()

	dispatch(c.Init())

	renderer := NewRenderer(w, h)

	for {
		c.Render(RenderContext{
			Viewport: NewViewport(renderer.Back),
			Focused:  true,
		})

		renderer.SwapBuffers()
		renderer.WriteTo(terminal)

		select {
		case <-ctx.Done():
			return nil
		case ev := <-ch:
			if ev == ShutdownEvent {
				cancel()
				return nil
			}

			task := c.Update(ev)
			dispatch(task)
		}
	}
}
