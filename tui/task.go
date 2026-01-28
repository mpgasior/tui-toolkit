package tui

import (
	"context"
)

type Task struct {
	ID      string
	Execute func(ctx context.Context, ch chan<- Event)
}

func None() Task {
	return Task{}
}

func Send(e Event) Task {
	return Task{
		Execute: func(ctx context.Context, ch chan<- Event) {
			select {
			case <-ctx.Done():
			case ch <- e:
			}
		},
	}
}
