package tui

import (
	"context"
)

type Task struct {
	ID      string
	Execute func(ctx context.Context, ch chan<- Event)
}

func TaskNone() Task {
	return Task{}
}

func TaskOne(e Event) Task {
	return Task{
		Execute: func(ctx context.Context, ch chan<- Event) {
			select {
			case <-ctx.Done():
			case ch <- e:
			}
		},
	}
}

func TaskF(f func(ctx context.Context, ch chan<- Event)) Task {
	return Task{
		Execute: f,
	}
}

var TaskShutdown = TaskOne(ShutdownEvent)
