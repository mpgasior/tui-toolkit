package ui

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/process"
	"github.com/mpgasior/tui-toolkit/mvu"
)

type ProcessStore struct {
	mu   sync.RWMutex
	data []process.ProcessInfo
}

type ProcessSnapshotEvent struct {
	ID   int
	Rows []process.ProcessInfo
}

type DataRefreshedEvent struct{}

func TaskFetchProcesses(store *ProcessStore) mvu.Task {
	return mvu.TaskNone
}

func TaskSearch(store *ProcessStore, term string) mvu.Task {
	return mvu.TaskNone
}

func TaskRefresh(id int, term []rune) mvu.Task {
	t := make([]rune, len(term))
	copy(t, term)
	return mvu.Task{
		ID: "search",
		Execute: func(ctx context.Context, ch chan<- mvu.Event) {
			for {
				info, err := process.List()
				if err != nil {
					continue
				}

				var filtered []process.ProcessInfo
				for _, p := range info {
					if len(t) == 0 || strings.Contains(p.Name, string(t)) {
						filtered = append(filtered, p)
					}
				}

				ev := ProcessSnapshotEvent{
					ID:   id,
					Rows: filtered,
				}

				select {
				case <-ctx.Done():
					return
				case ch <- ev:
				}

				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Second):
				}
			}
		},
	}
}
