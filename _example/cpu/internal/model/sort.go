package model

import (
	"cmp"
	"slices"
	"strings"
)

type SortOrder int

const (
	SortOrderAscending SortOrder = iota
	SortOrderDescending
)

func (s SortOrder) Next() SortOrder {
	return (s + 1) % 2
}

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

var sorters = map[SortBy]func(a, b Process) int{
	SortByPID:    func(a, b Process) int { return cmp.Compare(a.PID, b.PID) },
	SortByName:   func(a, b Process) int { return strings.Compare(a.Name, b.Name) },
	SortByAge:    func(a, b Process) int { return cmp.Compare(a.Age, b.Age) },
	SortByAvgCPU: func(a, b Process) int { return cmp.Compare(a.CPUAvg, b.CPUAvg) },
	SortByCPU:    func(a, b Process) int { return cmp.Compare(a.CPU, b.CPU) },
	SortByMem:    func(a, b Process) int { return cmp.Compare(a.MemoryRSS, b.MemoryRSS) },
	SortByMaxMem: func(a, b Process) int { return cmp.Compare(a.MemoryMax, b.MemoryMax) },
}

func SortResults(rows []Process, sortBy SortBy, order SortOrder) {
	fn, ok := sorters[sortBy]
	if !ok {
		return
	}

	slices.SortFunc(rows, sorters[SortByPID])

	slices.SortStableFunc(rows, func(a Process, b Process) int {
		result := fn(a, b)

		if order == SortOrderDescending {
			return -result
		}

		return result
	})
}
