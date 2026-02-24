package model

import (
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
)

type SortOrder int

const (
	SortOrderAscending SortOrder = iota
	SortOrderDescending
)

type SortBy int

const (
	SortByPID SortBy = iota
	SortByCPU
	SortByAvgCPU
	SortByName
	SortByAge
	SortByMem
	SortByMaxMem
)

type ListQuery struct {
	Term  string
	By    SortBy
	Order SortOrder
}

type Process struct {
	process.Snapshot

	AgeReady bool
	Age      time.Duration
}

type ProcessHistory struct {
	process.HistorySnapshot
}
