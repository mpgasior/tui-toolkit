package mvu

import (
	"context"

	"github.com/mpgasior/tui-toolkit/session"
	"github.com/mpgasior/tui-toolkit/termx"
	"github.com/mpgasior/tui-toolkit/vt"
)

type Session struct {
	terminal *termx.Terminal
	runtime  *Runtime

	dispatch func(Task)
	sess     *session.Session
}

func NewSession(term *termx.Terminal, ch chan<- Event) *Session {
	session := &Session{
		terminal: term,
		runtime:  NewRuntime(ch),
	}

	return session
}

func (s *Session) Start(ctx context.Context) error {
	dispatch, shutdown := s.runtime.Start(ctx)

	s.sess = session.New(
		func() (func() error, error) { return s.terminal.MakeRaw() },
		func() (func() error, error) { return vt.EnterMode(s.terminal, vt.ModeAlternateScreen) },
		func() (func() error, error) { return vt.EnterMode(s.terminal, vt.ModeBracketedPaste) },
		func() (func() error, error) { return vt.ExitMode(s.terminal, vt.ModeShowCursor) },
		func() (func() error, error) { return func() error { shutdown(); return nil }, nil },
	)

	if err := s.sess.Start(); err != nil {
		s.sess = nil
		shutdown()
		return err
	}

	eventF, eventCh := vt.Events(s.terminal)

	dispatch(TaskF(func(ctx context.Context, ch chan<- Event) {
		eventF(ctx)
	}))

	dispatch(TaskF(func(ctx context.Context, ch chan<- Event) {
		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-eventCh:
				select {
				case <-ctx.Done():
					return
				case ch <- ev:
				}
			}
		}
	}))
	dispatch(TaskF(func(ctx context.Context, ch chan<- Event) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-s.terminal.ResizeC():
				w, h, err := s.terminal.GetSize()
				if err != nil {
					continue
				}
				ev := ResizeEvent{
					Width: w, Height: h,
				}
				select {
				case <-ctx.Done():
					return
				case ch <- ev:
				}
			}
		}
	}))

	s.dispatch = dispatch

	return nil
}

func (s *Session) Stop() error {
	session := s.sess
	if session != nil {
		s.dispatch = nil
		s.sess = nil
		return session.Stop()
	}

	return nil
}

func (s *Session) RunSuspended(ctx context.Context, run func() error) (runErr error, err error) {
	if err = s.Stop(); err != nil {
		return nil, err
	}

	runErr = run()

	err = s.Start(ctx)
	return runErr, err
}

func (s *Session) Dispatch(t Task) {
	dispatch := s.dispatch
	if dispatch != nil {
		dispatch(t)
	}
}
