package model

import "time"

type State struct {
	Rows []QueryResult

	IsLoading  bool
	SearchTerm string
	SortBy     SortBy
	SortOrder  SortOrder

	IsPaused   bool
	LastUpdate time.Time
}

func (s *State) CurrentQuery() Query {
	return Query{
		Term:      s.SearchTerm,
		SortBy:    s.SortBy,
		Direction: s.SortOrder,
	}
}

func (s *State) PauseRefresh(pause bool) {
	s.IsPaused = pause
}

func (s *State) Sync(rows []QueryResult) (synced bool) {
	if s.IsPaused {
		return false
	}
	s.Rows = rows
	s.LastUpdate = time.Now()
	return true
}
