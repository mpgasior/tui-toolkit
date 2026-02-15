package tasks

import (
	"context"
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
	"github.com/mpgasior/tui-toolkit/mvu"
)

type DataRefreshedEvent struct{}

func TaskRefresh(store *process.Store, interval time.Duration) mvu.Task {
	return mvu.Task{
		ID: "refresh",
		Execute: func(ctx context.Context, ch chan<- mvu.Event) {
			for {
				snapshot, err := process.GetAll()
				if err != nil {
					continue
				}

				store.Sync(snapshot)

				select {
				case <-ctx.Done():
					return
				case ch <- DataRefreshedEvent{}:
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
