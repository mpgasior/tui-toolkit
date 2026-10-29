package termio

import (
	"context"
	"errors"
	"os"

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

func OpenSession(in, out *os.File) (*session, error) {
	s := &session{
		in: in, out: out,
	}

	if err := s.Enter(); err != nil {
		return nil, err
	}

	return s, nil
}

type session struct {
	in, out *os.File
	state   *term.State

	reader Reader
}

func (s *session) Enter() error {
	fd := int(s.in.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		return nil
	}
	s.state = state

	r, err := NewReader(fd)
	if err != nil {
		return errors.Join(
			err,
			term.Restore(fd, s.state))
	}
	s.reader = r

	return nil
}

func (s *session) Close() error {
	fd := int(s.in.Fd())
	return errors.Join(
		term.Restore(fd, s.state),
		s.reader.Close())
}

func (s *session) Suspend(fn func()) error {
	if err := s.Close(); err != nil {
		return err
	}

	fn()

	return s.Enter()
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
