package task

import (
	"context"
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/model"
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
	"github.com/mpgasior/tui-toolkit/mvu"
)

type HistoryReadyEvent struct {
	Key   process.Key
	Found bool
	Data  model.ProcessHistory
}

func QueryHistory(registry *process.Registry, key process.Key) mvu.Task {
	return mvu.Task{
		ID: "query-history",
		Execute: func(ctx context.Context, ch chan<- mvu.Event) {
			snapshot, found := registry.GetHistory(key)

			e := HistoryReadyEvent{
				Key:   key,
				Found: found,
				Data:  model.ProcessHistory{snapshot},
			}

			select {
			case <-ctx.Done():
				return
			case ch <- e:
			}
		},
	}
}
