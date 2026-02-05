package session

import "errors"

type SessionFn func() (restore func() error, err error)
type Session struct {
	steps   []SessionFn
	restore func() error
}

func New(funcs ...SessionFn) *Session {
	s := &Session{}

	for _, f := range funcs {
		s.steps = append(s.steps, f)
	}

	return s
}

func (s *Session) Start() error {
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

	for _, step := range s.steps {
		restore, err := step()
		if err != nil {
			return errors.Join(rollback(), err)
		}
		cleanupSteps = append(cleanupSteps, restore)
	}

	s.restore = rollback

	return nil
}

func (s *Session) Stop() error {
	restore := s.restore
	if restore != nil {
		s.restore = nil
		return restore()
	}

	return nil
}

func (s *Session) RunSuspended(run func() error) (runErr error, err error) {
	if err = s.Stop(); err != nil {
		return nil, err
	}

	runErr = run()

	err = s.Start()
	return runErr, err
}
