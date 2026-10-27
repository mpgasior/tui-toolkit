package async

import (
	"context"
	"sync"
)

type Runner struct {
	mu    sync.Mutex
	tasks map[FuncId]*contextImpl

	ctx context.Context

	out chan<- any
	wg  sync.WaitGroup
}

func NewRunner(ctx context.Context, out chan<- any) *Runner {
	return &Runner{
		ctx:   ctx,
		out:   out,
		tasks: make(map[FuncId]*contextImpl),
	}
}

func (r *Runner) Run(id FuncId, fn Func) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cancel(id)

	ctx, cancel := context.WithCancel(r.ctx)
	tCtx := &contextImpl{
		Context: ctx,
		cancel:  cancel,
		id:      id,
		runner:  r,
	}
	r.tasks[id] = tCtx

	r.wg.Go(func() {
		defer func() {
			r.mu.Lock()
			defer r.mu.Unlock()

			tCtx.cancel()
			if ctx, ok := r.tasks[id]; ok && ctx == tCtx {
				delete(r.tasks, id)
			}
		}()

		v, err := fn(tCtx)
		e := FuncCompletedEvent{
			Id:   id,
			Data: v,
			Err:  err,
		}
		r.yieldEvent(e)
	})
}

func (r *Runner) cancel(id FuncId) {
	if ctx, ok := r.tasks[id]; ok {
		ctx.cancel()
		delete(r.tasks, id)
	}
}

func (r *Runner) Cancel(id FuncId) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cancel(id)
}

func (r *Runner) CancelAll() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, ctx := range r.tasks {
		ctx.cancel()
		delete(r.tasks, id)
	}
}

func (r *Runner) yieldEvent(e any) {
	select {
	case <-r.ctx.Done():
		return
	case r.out <- e:
	}
}

func (r *Runner) Wait() {
	r.wg.Wait()
}
