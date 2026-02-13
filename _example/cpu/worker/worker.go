package worker

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/process"
	"github.com/mpgasior/tui-toolkit/mvu"
)

type DataRefreshedEvent struct{}

type QueryResultEvent struct {
	Rows []process.Profile
}

func TaskQuery(store *process.Store, term string) mvu.Task {
	execute := func(ctx context.Context, ch chan<- mvu.Event) {
		select {
		case <-ctx.Done():
			return
		case <-time.After(80 * time.Millisecond):
		}

		var filtered []process.Profile

		term = strings.ToLower(term)
		profiles := store.GetAll()

		for _, p := range profiles {
			name := strings.ToLower(p.Info.Name)
			if term == "" || strings.Contains(name, term) {
				filtered = append(filtered, p)
			}
		}

		sort.Slice(filtered, func(i, j int) bool {
			left, right := filtered[i], filtered[j]

			return left.History.RecentCPU() > right.History.RecentCPU()
		})

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
