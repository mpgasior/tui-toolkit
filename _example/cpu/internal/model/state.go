package model

import "time"

type State struct {
	Rows []QueryResult

	IsLoading  bool
	SearchTerm string
	SortBy     SortBy

	IsPaused   bool
	LastUpdate time.Time
}
