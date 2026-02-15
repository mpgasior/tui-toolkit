package model

import (
	"cmp"
	"slices"
	"strings"
)

var sorters = map[SortBy]func(a, b QueryResult) int{
	SortByPID:          func(a, b QueryResult) int { return cmp.Compare(a.PID, b.PID) },
	SortByName:         func(a, b QueryResult) int { return strings.Compare(a.Name, b.Name) },
	SortByCreationTime: func(a, b QueryResult) int { return a.CreationTime.Compare(a.CreationTime) },
	SortByAvgCPU:       func(a, b QueryResult) int { return cmp.Compare(a.AvgCPU, b.AvgCPU) },
	SortByRecentCPU:    func(a, b QueryResult) int { return cmp.Compare(a.RecentCPU, b.RecentCPU) },
}

func SortResults(rows []QueryResult, sortBy SortBy, order SortOrder) {
	fn, ok := sorters[sortBy]
	if !ok {
		return
	}

	slices.SortFunc(rows, func(a QueryResult, b QueryResult) int {
		result := fn(a, b)

		if order == SortOrderDescending {
			return -result
		}

		return result
	})
}
