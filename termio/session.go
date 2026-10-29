package termio

import (
	"context"
	"os"

	"github.com/mpgasior/tui-toolkit/scope"
	"golang.org/x/term"
)

type Session interface {
	In() *os.File
	Out() *os.File
	ReadContext(ctx context.Context, p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	GetSize() (w, h int, err error)
	Close() error
}

type SessionModifier struct {
	Apply   func(s Session) error
	Restore func(s Session) error
}

func OpenSession(in, out *os.File, mods ...SessionModifier) (*session, error) {
	s := &session{
		in: in, out: out,
		mods: mods,
	}

	if err := s.enter(); err != nil {
		return nil, err
	}

	return s, nil
}

type session struct {
	in, out *os.File
	state   *term.State
	mods    []SessionModifier
	scope   *scope.Scope[*session]

	reader Reader
}

func (s *session) enter() error {
	builder := scope.NewBuilder(s,
		scope.Step[*session]{
			Apply: func(s *session) error {
				fd := int(s.in.Fd())
				state, err := term.MakeRaw(fd)
				if err != nil {
					return err
				}
				s.state = state
				return nil
			},
			Restore: func(s *session) error {
				fd := int(s.in.Fd())
				return term.Restore(fd, s.state)
			},
		},
		scope.Step[*session]{
			Apply: func(s *session) error {
				fd := int(s.in.Fd())
				var err error
				s.reader, err = NewReader(fd)
				return err
			},
			Restore: func(s *session) error {
				return s.reader.Close()
			},
		})

	toStep := func(fn func(Session) error) func(*session) error {
		return func(s *session) error { return fn(s) }
	}

	for _, mod := range s.mods {
		builder.AddStep(scope.Step[*session]{
			Apply:   toStep(mod.Apply),
			Restore: toStep(mod.Restore),
		})
	}

	scope := builder.Scope()
	if err := scope.Apply(); err != nil {
		return err
	}

	s.scope = scope
	return nil
}

func (s *session) Close() error {
	return s.scope.Restore()
}

func (s *session) Suspend(fn func()) error {
	if err := s.Close(); err != nil {
		return err
	}

	fn()

	return s.enter()
}

func (s *session) In() *os.File {
	return s.in
}

func (s *session) Out() *os.File {
	return s.out
}

func (s *session) ReadContext(ctx context.Context, p []byte) (n int, err error) {
	return s.reader.ReadContext(ctx, p)
}

func (s *session) Write(p []byte) (n int, err error) {
	return s.out.Write(p)
}

func (s *session) GetSize() (w, h int, err error) {
	fd := int(s.out.Fd())
	return term.GetSize(fd)
}
