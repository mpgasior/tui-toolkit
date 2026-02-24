package task

import (
	"context"
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
	"github.com/mpgasior/tui-toolkit/mvu"
)

type RegistryRefreshedEvent struct{}

func CancelRefresh() mvu.Task {
	return mvu.TaskCancel("refresh")
}

func Refresh(registry *process.Registry, interval time.Duration) mvu.Task {
	return mvu.Task{
		ID: "refresh",
		Execute: func(ctx context.Context, ch chan<- mvu.Event) {
			for {
				samples, err := process.GetAll()
				if err != nil {
					continue
				}

				registry.Update(samples)

				select {
				case <-ctx.Done():
					return
				case ch <- RegistryRefreshedEvent{}:
				}

				select {
				case <-ctx.Done():
					return
				case <-time.After(interval):
				}
			}
		},
	}
}
