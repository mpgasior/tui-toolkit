package tui

import (
	"context"
)

type Task struct {
	ID      string
	Execute func(ctx context.Context, ch chan<- Event)
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

var TaskNone = Task{}
var TaskShutdown = TaskOne(ShutdownEvent)
