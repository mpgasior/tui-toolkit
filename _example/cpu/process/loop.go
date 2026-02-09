package process

import (
	"context"
	"time"
)

func RefreshLoop(ctx context.Context, ch chan []ProcessInfo, idle time.Duration) {
	for {
		info, err := List()
		if err != nil {
			return
		}
		select {
		case <-ctx.Done():
			return
		case ch <- info:
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(idle):
		}
	}
}
