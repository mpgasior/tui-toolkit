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
	SortByRecentMemory
)

type Query struct {
	Term      string
	SortBy    SortBy
	Direction SortOrder
}

type QueryResult struct {
	PID          uint32
	Name         string
	CreationTime time.Time
	AvgCPU       float64
	RecentCPU    float64
}
