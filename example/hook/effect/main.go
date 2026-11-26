package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/mpgasior/tui-toolkit/ansi"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/event"
	"github.com/mpgasior/tui-toolkit/gfx"
	"github.com/mpgasior/tui-toolkit/hook"
	"github.com/mpgasior/tui-toolkit/seq"
	"github.com/mpgasior/tui-toolkit/termio"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	term, err := termio.OpenSession(os.Stdin, os.Stdout, termio.WithAltScreen)
	if err != nil {
		return err
	}
	defer term.Close()

	w, h, err := term.GetSize()
	if err != nil {
		return err
	}

	screen := gfx.NewScreen(w, h, nil)
	escReader := seq.NewEscReader(context.TODO(), term.ReadContext)
	scanner := bufio.NewScanner(escReader)
	scanner.Split(seq.Scan)
	keyDecoder := event.NewKeyDecoder()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	errChan := make(chan error)
	keyChan := make(chan *hook.KeyEvent)

	events := func() {
		for ctx.Err() == nil {
			if scanner.Scan() {
				token := scanner.Bytes()
				if key, ok := keyDecoder.Decode(token); ok {
					if key.Code == ansi.KeyCtrlC {
						cancel()
						return
					}
					event := &hook.KeyEvent{
						Code: key.Code,
						R:    key.R,
					}
					keyChan <- event
				}
			}

			if err := scanner.Err(); err != nil {
				errChan <- err
				return
			}
		}
	}

	renderChan := make(chan any, 1)
	root := hook.Root(func() {
		select {
		case renderChan <- true:
		default:
		}
	})
	defer root.Unmount()

	var wg sync.WaitGroup
	wg.Go(events)
	defer wg.Wait()

	for {
		sf := draw.NewSurface(screen.Back)
		render := root.Render(app)
		render(sf)

		screen.Swap()
		if _, err := screen.WriteTo(term); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return nil
		case err := <-errChan:
			return err
		case key := <-keyChan:
			hook.DispatchKey(root, key)
		case <-renderChan:
		}
	}
}

func app(s *hook.Scope) hook.RenderFunc {
	down, setDown := hook.State(s, false)

	counter := s.Mount("counter", counter(down, 100))
	clock := s.Mount("timer", clock)

	hook.OnKey(s, func(e *hook.KeyEvent) {
		switch e.Code {
		case ansi.KeyUp:
			hook.SetState(setDown, false)
		case ansi.KeyDown:
			hook.SetState(setDown, true)
		}
	})

	return func(sf *draw.Surface) {
		draw.Clear(sf)

		counter(sf)
		clockSlot := sf.Inset(1, 0, 0, 0)
		clock(clockSlot)
	}
}

func counter(down bool, init int) hook.Component {
	return func(s *hook.Scope) hook.RenderFunc {
		count, setCount := hook.State(s, init)

		hook.Effect(s, func() hook.EffectCleanup {
			ctx, cancel := context.WithCancel(context.Background())
			var wg sync.WaitGroup

			wg.Go(func() {
				timer := time.NewTicker(time.Second)
				defer timer.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-timer.C:
						setCount(func(prev int) int {
							delta := 1
							if down {
								delta = -1
							}

							return prev + delta
						})
					}
				}
			})
			return func() {
				cancel()
				wg.Wait()
			}
		}, down)

		return func(sf *draw.Surface) {
			w := draw.NewTextWriter(sf)
			fmt.Fprintf(w, "%d", count)
		}
	}
}

func clock(s *hook.Scope) hook.RenderFunc {
	now, setNow := hook.StateFn(s, time.Now)
	hook.Effect(s, func() hook.EffectCleanup {
		timer := time.NewTicker(time.Second)
		done := make(chan any)
		var wg sync.WaitGroup

		wg.Go(func() {
			for {
				select {
				case <-timer.C:
					setNow(func(_ time.Time) time.Time { return time.Now() })
				case <-done:
					return
				}
			}
		})

		return func() {
			close(done)
			timer.Stop()
			wg.Wait()
		}
	})

	return func(sf *draw.Surface) {
		w := draw.NewTextWriter(sf)
		fmt.Fprintf(w, "%v", now)
	}
}
