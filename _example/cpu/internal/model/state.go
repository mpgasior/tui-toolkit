package model

import (
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
)

type State struct {
	Store *process.Store

	PID uint32

	IsPaused   bool
	Filtered   []ProcessSummary
	LastUpdate time.Time

	SearchTerm string
	SortBy     SortBy
	SortOrder  SortOrder
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

func (s *State) Sync(results []ProcessSummary) {
	s.Filtered = results
	s.LastUpdate = time.Now()
}
