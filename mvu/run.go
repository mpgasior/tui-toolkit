package mvu

import (
	"bytes"
	"context"

	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/termx"
	"github.com/mpgasior/tui-toolkit/view"
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

	scr := screen.New(w, h)

	for {
		c.Render(RenderContext{
			View:    view.NewPort(scr),
			Focused: true,
		})

		scr.WriteTo(terminal)

		select {
		case <-ctx.Done():
			return nil
		case ev := <-ch:
			if ev == ShutdownEvent {
				cancel()
				return nil
			}

			if resize, ok := ev.(ResizeEvent); ok {
				scr = screen.New(resize.Width, resize.Height)
				scr.Flush(terminal)
				continue
			}

			if batch, ok := ev.(BatchTaskEvent); ok {
				for _, t := range batch.Tasks {
					session.Dispatch(t)
				}
				continue
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

				if _, err := scr.Flush(terminal); err != nil {
					return err
				}
				continue
			}

			session.Dispatch(c.Update(ev))
		}
	}
}
