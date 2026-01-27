package tui

import "context"

type Task func(ctx context.Context, ch chan<- Event)

func None() Task {
	return nil
}

func Send(e Event) Task {
	return func(ctx context.Context, ch chan<- Event) {
		select {
		case <-ctx.Done():
		case ch <- e:
		}
	}
}
