package tui

import (
	"bytes"
	"context"

	"github.com/mpgasior/tui-go/termx"
)

func Run(c Component) error {
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

	ch := make(chan Event)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session := NewSession(terminal, ch)
	if err := session.Start(ctx); err != nil {
		return err
	}
	defer session.Stop()

	session.Dispatch(c.Init())

	w, h, err := terminal.GetSize()
	if err != nil {
		return err
	}

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
				var buf bytes.Buffer
				runErr, err := session.RunSuspended(ctx, func() error {
					cmd, captureOutput, err := launch.CmdBuilder(tty.In, tty.Out)
					if err != nil {
						return err
					}

					if captureOutput {
						cmd.Stdout = &buf
					}

					return cmd.Run()
				})

				if err != nil {
					return err
				}

				if launch.OnResult != nil {
					t := launch.OnResult(buf.Bytes(), runErr)
					session.Dispatch(t)
				}

				continue
			}

			session.Dispatch(c.Update(ev))
		}
	}
}
