package hook

import "context"

func Go(s *Scope, fn func(ctx context.Context)) context.CancelFunc {
	ctx, cancel := context.WithCancel(s.context)
	s.wg.Go(func() {
		fn(ctx)
	})

	return cancel
}
