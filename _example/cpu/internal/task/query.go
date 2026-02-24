package task

import (
	"context"
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/model"
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
	"github.com/mpgasior/tui-toolkit/mvu"
)

type ListReadyEvent struct {
	Query model.ProcessListQuery
	Data  []model.Process
}

func RebuildSnapshot(registry *process.Registry, query model.ProcessListQuery) mvu.Task {
	return mvu.Task{
		ID: "query",
		Execute: func(ctx context.Context, ch chan<- mvu.Event) {
			select {
			case <-ctx.Done():
				return
			case <-time.After(80 * time.Millisecond):
			}

			snapshot := registry.GetSnapshot()
			results := make([]model.Process, 0, len(snapshot))

			for _, s := range snapshot {
				results = append(results, toProcess(s))
			}

			results = model.Filter(results, query.Term)
			model.SortResults(results, query.By, query.Order)

			e := ListReadyEvent{
				Query: query,
				Data:  results,
			}

			select {
			case <-ctx.Done():
				return
			case ch <- e:
			}
		},
	}
}

func toProcess(s process.Snapshot) model.Process {
	ageReady := false
	age := time.Duration(0)
	if !s.CreationTime.IsZero() {
		ageReady = true
		age = time.Since(s.CreationTime)
	}

	return model.Process{
		Snapshot: s,
		Age:      age,
		AgeReady: ageReady,
	}
}
