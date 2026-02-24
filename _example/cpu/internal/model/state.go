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

func (s *State) TogglePause() bool {
	s.IsPaused = !s.IsPaused
	return s.IsPaused
}

func (s *State) CurrentListQuery() ProcessListQuery {
	return ProcessListQuery{
		Term:  s.SearchTerm,
		By:    s.SortBy,
		Order: s.SortOrder,
	}
}
