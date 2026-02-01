package tui

import (
	"context"
	"errors"

	"github.com/mpgasior/tui-go/termx"
	"github.com/mpgasior/tui-go/vt"
)

type Session struct {
	terminal *termx.Terminal
	runtime  *Runtime

	dispatch func(Task)
	restore  func() error
}

func NewSession(term *termx.Terminal, ch chan<- Event) *Session {
	session := &Session{
		terminal: term,
		runtime:  NewRuntime(ch),
	}

	return session
}

func (s *Session) Start(ctx context.Context) error {
	var cleanupSteps []func() error
	rollback := func() error {
		var errs []error
		for i := len(cleanupSteps) - 1; i >= 0; i-- {
			if err := cleanupSteps[i](); err != nil {
				errs = append(errs, err)
			}
		}
		return errors.Join(errs...)
	}

	steps := []struct {
		action func() (func() error, error)
	}{
		{action: func() (func() error, error) { return s.terminal.MakeRaw() }},
		{action: func() (func() error, error) { return vt.EnterMode(s.terminal, vt.ModeAlternateScreen) }},
		{action: func() (func() error, error) { return vt.EnterMode(s.terminal, vt.ModeBracketedPaste) }},
		{action: func() (func() error, error) { return vt.ExitMode(s.terminal, vt.ModeShowCursor) }},
	}

	for _, step := range steps {
		restore, err := step.action()
		if err != nil {
			return errors.Join(rollback(), err)
		}
		cleanupSteps = append(cleanupSteps, restore)
	}

	dispatch, shutdown := s.runtime.Start(ctx)
	cleanupSteps = append(cleanupSteps, func() error { shutdown(); return nil })
	dispatch(TaskF(Input(s.terminal)))

	s.dispatch = dispatch
	s.restore = rollback

	return nil
}

func (s *Session) Stop() error {
	restore := s.restore
	if restore != nil {
		s.restore = nil
		s.dispatch = nil
		return restore()
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
