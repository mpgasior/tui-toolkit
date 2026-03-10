package model

import (
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
)

type State struct {
	Registry *process.Registry

	IsPaused bool

	SearchTerm string
	SortBy     SortBy
	SortOrder  SortOrder
	Exclude    Exclude

	SelectedKey process.Key
}

func New(size int, recycleAfter time.Duration) State {
	return State{
		Registry: process.NewRegistry(size, recycleAfter),
	}
}

func (s *State) UpdateSort(by SortBy, order SortOrder, exclude Exclude) {
	s.SortBy = by
	s.SortOrder = order
	s.Exclude = exclude
}

func (s *State) TogglePause() bool {
	s.IsPaused = !s.IsPaused
	return s.IsPaused
}

func (s *State) UpdateExclude(exclude Exclude) {
	s.Exclude = exclude
}

func (s *State) CurrentListQuery() ListQuery {
	return ListQuery{
		Term:    s.SearchTerm,
		By:      s.SortBy,
		Order:   s.SortOrder,
		Exclude: s.Exclude,
	}
}
