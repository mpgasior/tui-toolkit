package model

import "time"

type SortOrder int

const (
	SortOrderAscending SortOrder = iota
	SortOrderDescending
)

type SortBy int

const (
	SortByPID SortBy = iota
	SortByRecentCPU
	SortByAvgCPU
	SortByName
	SortByAge
	SortByWorkingSet
	SortByPeakWorkingSet
)

type ProcessListQuery struct {
	Term  string
	By    SortBy
	Order SortOrder
}

type ProcessSummary struct {
	PID  uint32
	Name string
	Age  time.Duration

	Computed       bool
	AvgCPU         float64
	RecentCPU      float64
	WorkingSet     uint64
	PeakWorkingSet uint64
}
