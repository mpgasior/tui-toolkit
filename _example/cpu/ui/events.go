package ui

import (
	"context"
	"strings"
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/process"
	"github.com/mpgasior/tui-toolkit/mvu"
)

type ProcessSnapshotEvent struct {
	ID   int
	Rows []process.ProcessInfo
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
