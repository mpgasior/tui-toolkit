package mvu

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

func TaskCancel(id string) Task {
	return Task{ID: id}
}

var TaskNone = Task{}
var TaskShutdown = TaskOne(ShutdownEvent)
