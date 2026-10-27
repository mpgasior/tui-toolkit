package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/mpgasior/tui-toolkit/async"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ch := make(chan any)
	r := async.NewRunner(ctx, ch)

	tick := func(period time.Duration) func(ctx async.Context) (any, error) {
		return func(ctx async.Context) (any, error) {
			t := time.NewTicker(period)
			defer t.Stop()

			c := 1
			for {
				select {
				case <-ctx.Done():
					return nil, nil
				case <-t.C:
					ctx.YieldEvent(fmt.Sprintf("%s: %d", ctx.Id(), c))
					c += 1
				}
			}

			return nil, nil
		}
	}

	r.Run("second", tick(time.Second))
	r.Run("half-second", tick(time.Millisecond*500))

	go func() {
		for e := range ch {
			fmt.Println(e)
		}
	}()

	context.AfterFunc(ctx, func() {
		r.CancelAll()
	})

	r.Wait()
}
