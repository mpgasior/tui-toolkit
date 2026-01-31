package tui

import (
	"bytes"
	"context"
	"os"

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
	defer terminal.Close()

	restore, err := terminal.MakeRaw()
	if err != nil {
		return err
	}
	defer restore()

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
	dispatch(TaskF(Input(terminal)))

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

			if launch, ok := ev.(LaunchEvent); ok {
				shutdown()
				restore()

				cmd, captureOutput, _ := launch.CmdBuilder()
				if cmd.Stderr == nil {
					cmd.Stderr = os.Stderr
				}
				if cmd.Stdout == nil {
					cmd.Stdout = tty.Out
				}
				if cmd.Stdin == nil {
					cmd.Stdin = tty.In
				}

				var buf bytes.Buffer
				if captureOutput {
					cmd.Stdout = &buf
				}

				execErr := cmd.Run()

				dispatch, shutdown = runtime.Start(ctx)
				defer shutdown()

				restore, err = terminal.MakeRaw()
				if err != nil {
					return err
				}
				defer restore()
				dispatch(TaskF(Input(terminal)))

				if launch.OnResult != nil {
					t := launch.OnResult(buf.Bytes(), execErr)
					dispatch(t)
				}

				continue
			}

			task := c.Update(ev)
			dispatch(task)
		}
	}
}
