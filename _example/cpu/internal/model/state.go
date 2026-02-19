package model

import (
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
)

type State struct {
	Store *process.Store

	PID uint32

	IsPaused   bool
	Filtered   []QueryResult
	LastUpdate time.Time

	SearchTerm string
	SortBy     SortBy
	SortOrder  SortOrder
}

func NewState() *State {
	return &State{
		Store: process.NewStore(),
	}
}

func (s *State) TogglePause() bool {
	s.IsPaused = !s.IsPaused
	return s.IsPaused
}

func (s *State) CurrentQuery() Query {
	return Query{
		Term:      s.SearchTerm,
		SortBy:    s.SortBy,
		Direction: s.SortOrder,
	}
}

func (s *State) Sync(results []QueryResult) {
	s.Filtered = results
	s.LastUpdate = time.Now()
}
