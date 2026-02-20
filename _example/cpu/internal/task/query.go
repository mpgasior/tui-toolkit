package task

import (
	"context"
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/model"
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
	"github.com/mpgasior/tui-toolkit/mvu"
)

type QuerySingleResultEvent struct {
	PID    uint32
	Found  bool
	Result process.Profile
}

func QuerySingle(store *process.Store, pid uint32) mvu.Task {
	return mvu.Task{
		ID: "query-single",
		Execute: func(ctx context.Context, ch chan<- mvu.Event) {
			profile, found := store.GetProfile(pid)

			ev := QuerySingleResultEvent{
				PID:    pid,
				Found:  found,
				Result: profile,
			}

			select {
			case <-ctx.Done():
				return
			case ch <- ev:
				return
			}
		},
	}
}

type ProcessSummaryEvent struct {
	Query model.ProcessListQuery
	Data  []model.ProcessSummary
}

func QueryProcessList(store *process.Store, query model.ProcessListQuery) mvu.Task {
	return mvu.Task{
		ID: "query",
		Execute: func(ctx context.Context, ch chan<- mvu.Event) {
			select {
			case <-ctx.Done():
				return
			case <-time.After(80 * time.Millisecond):
			}

			snapshot := store.GetAll()
			results := make([]model.ProcessSummary, 0, len(snapshot))

			for _, s := range snapshot {
				results = append(results, toProcessSummary(s))
			}

			results = model.Filter(results, query.Term)
			model.SortResults(results, query.By, query.Order)

			e := ProcessSummaryEvent{
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

func toProcessSummary(s process.Snapshot) model.ProcessSummary {
	age := time.Duration(0)
	if !s.CreationTime.IsZero() {
		age = time.Since(s.CreationTime)
	}

	return model.ProcessSummary{
		PID:  s.Info.PID,
		Name: s.Info.Name,
		Age:  age,

		Computed:       s.Computed,
		AvgCPU:         s.AvgCPU,
		RecentCPU:      s.RecentCPU,
		WorkingSet:     s.WorkingSet,
		PeakWorkingSet: s.PeakWorkingSet,
	}
}
