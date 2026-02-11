package worker

import (
	"context"
	"strings"
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/process"
	"github.com/mpgasior/tui-toolkit/mvu"
)

type DataRefreshedEvent struct{}

type QueryResultEvent struct {
	Rows []process.ProcessInfo
}

func TaskQuery(store *process.ProcessStore, term string) mvu.Task {
	execute := func(ctx context.Context, ch chan<- mvu.Event) {
		select {
		case <-ctx.Done():
			return
		case <-time.After(80 * time.Millisecond):
		}

		var filtered []process.ProcessInfo

		term = strings.ToLower(term)
		list := store.GetAll()

		for _, p := range list {
			name := strings.ToLower(p.Name)
			if term == "" || strings.Contains(name, term) {
				filtered = append(filtered, p)
			}
		}

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

func TaskRefresh(s *process.ProcessStore) mvu.Task {
	execute := func(ctx context.Context, ch chan<- mvu.Event) {
		for {
			list, err := process.GetAll()
			if err != nil {
				continue
			}

			s.Update(list)

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
