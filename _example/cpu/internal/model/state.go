package model

import (
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
)

type State struct {
	Registry *process.Registry

	IsPaused bool

	SearchTerm string
	SortBy     SortBy
	SortOrder  SortOrder

	SelectedKey process.Key
}

func New(size int) State {
	return State{
		Registry: process.NewRegistry(size),
	}
}

func (s *State) UpdateSort(by SortBy, order SortOrder) {
	s.SortBy = by
	s.SortOrder = order
}

func (s *State) TogglePause() bool {
	s.IsPaused = !s.IsPaused
	return s.IsPaused
}

func (s *State) CurrentListQuery() ListQuery {
	return ListQuery{
		Term:  s.SearchTerm,
		By:    s.SortBy,
		Order: s.SortOrder,
	}
}
