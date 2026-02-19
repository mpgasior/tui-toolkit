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
	Result model.QueryResult
}

func QuerySingle(store *process.Store, pid uint32) mvu.Task {
	return mvu.Task{
		ID: "query-single",
		Execute: func(ctx context.Context, ch chan<- mvu.Event) {
			snapshot := store.GetAll()

			for _, s := range snapshot {
				if s.Info.PID == pid {
					ev := QuerySingleResultEvent{
						PID:    pid,
						Result: toQueryResult(s),
					}
					select {
					case <-ctx.Done():
						return
					case ch <- ev:
					}
				}
			}
		},
	}
}

type QueryResultEvent struct {
	Query model.Query
	Data  []model.QueryResult
}

func Query(store *process.Store, query model.Query) mvu.Task {
	return mvu.Task{
		ID: "query",
		Execute: func(ctx context.Context, ch chan<- mvu.Event) {
			select {
			case <-ctx.Done():
				return
			case <-time.After(80 * time.Millisecond):
			}

			snapshot := store.GetAll()
			results := make([]model.QueryResult, 0, len(snapshot))

			for _, s := range snapshot {
				results = append(results, toQueryResult(s))
			}

			results = model.Filter(results, query.Term)
			model.SortResults(results, query.SortBy, query.Direction)

			e := QueryResultEvent{
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

func toQueryResult(s process.Snapshot) model.QueryResult {
	age := time.Duration(0)
	if !s.Info.CreationTime.IsZero() {
		age = time.Since(s.Info.CreationTime)
	}

	return model.QueryResult{
		PID:  s.Info.PID,
		Name: s.Info.Name,
		Age:  age,

		IsReady:        s.IsReady,
		AvgCPU:         s.AvgCPU,
		RecentCPU:      s.RecentCPU,
		WorkingSet:     s.WorkingSet,
		PeakWorkingSet: s.PeakWorkingSet,
	}
}
