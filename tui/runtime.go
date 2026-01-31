package tui

import (
	"context"
	"sync"
)

type Runtime struct {
	events  chan<- Event
	mu      sync.Mutex
	tasks   map[string]context.CancelFunc
	pending map[string]Task
}

func NewRuntime(ch chan<- Event) *Runtime {
	return &Runtime{
		tasks:   make(map[string]context.CancelFunc),
		pending: make(map[string]Task),
		events:  ch,
	}
}

func (r *Runtime) Start(ctx context.Context) (dispatch func(Task), shutdown func()) {
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	quitCh := make(chan string)

	run := func(t Task) {
		if t.ID == "" {
			wg.Go(func() { t.Execute(ctx, r.events) })
			return
		}

		tCtx, tCancel := context.WithCancel(ctx)
		r.tasks[t.ID] = tCancel

		wg.Go(func() {
			defer tCancel()

			t.Execute(tCtx, r.events)
			select {
			case quitCh <- t.ID:
			case <-tCtx.Done():
			}
		})
	}

	coordinate := func() {
		for {
			select {
			case <-ctx.Done():
				return
			case id := <-quitCh:
				r.mu.Lock()

				delete(r.tasks, id)

				if t, ok := r.pending[id]; ok {
					delete(r.pending, id)
					run(t)
				}

				r.mu.Unlock()
			}
		}
	}

	dispatch = func(t Task) {
		r.mu.Lock()
		defer r.mu.Unlock()

		if cancel, running := r.tasks[t.ID]; running {
			cancel()
			if t.Execute != nil {
				r.pending[t.ID] = t
			} else {
				delete(r.pending, t.ID)
			}
		} else if t.Execute != nil {
			run(t)
		}
	}

	shutdown = sync.OnceFunc(func() {
		cancel()
		wg.Wait()
	})

	wg.Go(coordinate)
	return dispatch, shutdown
}
