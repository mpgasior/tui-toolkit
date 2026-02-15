package tasks

import (
	"context"
	"time"

	"github.com/mpgasior/tui-toolkit/mvu"
)

type TickEvent struct {
	ID string
}

func TaskTick(id string, interval time.Duration) mvu.Task {
	return mvu.Task{
		ID: id,
		Execute: func(ctx context.Context, ch chan<- mvu.Event) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(interval):
				}

				select {
				case <-ctx.Done():
					return
				case ch <- TickEvent{ID: id}:
				}
			}
		},
	}
}
