package worker

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/process"
	"github.com/mpgasior/tui-toolkit/mvu"
)

type Query struct {
	Term   string
	SortBy string
}

type QueryResult struct {
	PID          uint32
	Name         string
	CreationTime time.Time
	AvgCPU       float64
	RecentCPU    float64
}

var sorters = map[string]func(a, b QueryResult) int{
	"PID":          func(a, b QueryResult) int { return cmp.Compare(b.PID, a.PID) },
	"Name":         func(a, b QueryResult) int { return strings.Compare(b.Name, a.Name) },
	"CreationTime": func(a, b QueryResult) int { return b.CreationTime.Compare(a.CreationTime) },
	"AvgCPU":       func(a, b QueryResult) int { return cmp.Compare(b.AvgCPU, a.AvgCPU) },
	"RecentCPU":    func(a, b QueryResult) int { return cmp.Compare(b.RecentCPU, a.RecentCPU) },
}

type DataRefreshedEvent struct{}

type QueryResultEvent struct {
	Rows []QueryResult
}

func TaskQuery(store *process.Store, query Query) mvu.Task {
	execute := func(ctx context.Context, ch chan<- mvu.Event) {
		select {
		case <-ctx.Done():
			return
		case <-time.After(80 * time.Millisecond):
		}

		var filtered []QueryResult

		term := strings.ToLower(query.Term)
		profiles := store.GetAll()

		for _, p := range profiles {
			name := strings.ToLower(p.Info.Name)
			if query.Term == "" || strings.Contains(name, term) {
				stats := p.History.Stats()
				filtered = append(filtered, QueryResult{
					PID:          p.Info.PID,
					Name:         p.Info.Name,
					CreationTime: p.Info.CreationTime,
					AvgCPU:       stats.AvgCPU,
					RecentCPU:    stats.RecentCPU,
				})
			}
		}

		slices.SortFunc(filtered, sorters[query.SortBy])

		ev := QueryResultEvent{
			Rows: filtered,
		}

		select {
		case <-ctx.Done():
			return
		case ch <- ev:
		}
	}

	return mvu.Task{
		ID:      "refresh",
		Execute: execute,
	}
}

func TaskRefresh(s *process.Store) mvu.Task {
	execute := func(ctx context.Context, ch chan<- mvu.Event) {
		for {
			snapshot, err := process.GetAll()
			if err != nil {
				continue
			}

			s.Sync(snapshot)

			select {
			case <-ctx.Done():
				return
			case ch <- DataRefreshedEvent{}:
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
		}
	}

	return mvu.Task{
		ID:      "background-refresh",
		Execute: execute,
	}
}
