package tui

import "context"

type Task func(ctx context.Context, ch chan<- Event)
