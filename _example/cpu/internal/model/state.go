package model

type State struct {
	Rows       []QueryResult
	IsLoading  bool
	SearchTerm string
	SortBy     SortBy
}
