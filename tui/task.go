package tui

import (
	"context"
	"sync"
)

type Task struct {
	ID      string
	Execute func(ctx context.Context, ch chan<- Event)
}

func None() Task {
	return Task{}
}

func Send(e Event) Task {
	return Task{
		Execute: func(ctx context.Context, ch chan<- Event) {
			select {
			case <-ctx.Done():
			case ch <- e:
			}
		},
	}
}

type Runtime struct {
	tasks   map[string]context.CancelFunc
	pending map[string]Task
	quitCh  chan string
	ch      chan<- Event
	ctx     context.Context
	cancel  context.CancelFunc
	mu      sync.Mutex
	wg      sync.WaitGroup
}

func NewRuntime(ctx context.Context, ch chan<- Event) *Runtime {
	ctx, cancel := context.WithCancel(ctx)
	return &Runtime{
		tasks:   make(map[string]context.CancelFunc),
		pending: make(map[string]Task),
		quitCh:  make(chan string),
		ch:      ch,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (r *Runtime) Start() {
	r.wg.Go(func() {
		for {
			select {
			case <-r.ctx.Done():
				return
			case id := <-r.quitCh:
				r.mu.Lock()

				delete(r.tasks, id)

				if t, ok := r.pending[id]; ok {
					delete(r.pending, id)
					r.startTask(t)
				}

				r.mu.Unlock()
			}
		}
	})
}

func (r *Runtime) Dispatch(t Task) {
	if t.ID == "" {
		if t.Execute != nil {
			r.wg.Go(func() {
				t.Execute(r.ctx, r.ch)
			})
		}
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	cancel, running := r.tasks[t.ID]

	if running {
		cancel()
		if t.Execute != nil {
			r.pending[t.ID] = t
		} else {
			delete(r.pending, t.ID)
		}
	} else if t.Execute != nil {
		r.startTask(t)
	}
}

func (r *Runtime) startTask(t Task) {
	ctx, cancel := context.WithCancel(r.ctx)
	r.tasks[t.ID] = cancel

	r.wg.Go(func() {
		defer func() {
			select {
			case r.quitCh <- t.ID:
			case <-r.ctx.Done():
			}
		}()

		t.Execute(ctx, r.ch)
	})
}

func (r *Runtime) Shutdown() {
	r.cancel()
	r.wg.Wait()
}
