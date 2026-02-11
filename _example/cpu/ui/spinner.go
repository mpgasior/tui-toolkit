package ui

import (
	"context"
	"time"

	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/mvu"
	"github.com/mpgasior/tui-toolkit/screen"
)

type SpinnerTickEvent struct {
	ID string
}

type Spinner struct {
	ID         string
	Frame      draw.SpinnerFrame
	FrameIndex int
}

func (s *Spinner) OnTick(e SpinnerTickEvent) {
	if e.ID == s.ID {
		s.FrameIndex += 1
	}
}

func (s *Spinner) Render(m screen.Mutator) {
	draw.Spinner(m, s.FrameIndex, s.Frame, screen.DefaultStyle)
}

func (s *Spinner) StartTask() mvu.Task {
	return mvu.Task{
		ID: s.ID,
		Execute: func(ctx context.Context, ch chan<- mvu.Event) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(80 * time.Millisecond):
				}

				select {
				case <-ctx.Done():
					return
				case ch <- SpinnerTickEvent{ID: s.ID}:
				}
			}
		},
	}
}

func (s *Spinner) CancelTask() mvu.Task {
	return mvu.TaskCancel(s.ID)
}
