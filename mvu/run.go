package mvu

import (
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
		case e := <-ch:
			switch ev := e.(type) {
			case shutdownEvent:
				cancel()
				return nil
			case ResizeEvent:
				scr = screen.New(ev.Width, ev.Height)
				scr.Flush(terminal)
			case BatchTaskEvent:
				for _, t := range ev.Tasks {
					session.Dispatch(t)
				}
			case ExecEvent:
				var task Task
				err := session.RunSuspended(ctx, func() {
					task = ev(tty.In, tty.Out)
				})
				if err != nil {
					return err
				}

				session.Dispatch(task)

				if _, err := scr.Flush(terminal); err != nil {
					return err
				}
			default:
				session.Dispatch(c.Update(ev))
			}
		}
	}
}
