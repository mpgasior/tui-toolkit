package model

import (
	"cmp"
	"slices"
	"strings"
)

var sorters = map[SortBy]func(a, b ProcessSummary) int{
	SortByPID:            func(a, b ProcessSummary) int { return cmp.Compare(a.PID, b.PID) },
	SortByName:           func(a, b ProcessSummary) int { return strings.Compare(a.Name, b.Name) },
	SortByAge:            func(a, b ProcessSummary) int { return cmp.Compare(a.Age, b.Age) },
	SortByAvgCPU:         func(a, b ProcessSummary) int { return cmp.Compare(a.AvgCPU, b.AvgCPU) },
	SortByRecentCPU:      func(a, b ProcessSummary) int { return cmp.Compare(a.RecentCPU, b.RecentCPU) },
	SortByWorkingSet:     func(a, b ProcessSummary) int { return cmp.Compare(a.WorkingSet, b.WorkingSet) },
	SortByPeakWorkingSet: func(a, b ProcessSummary) int { return cmp.Compare(a.PeakWorkingSet, b.PeakWorkingSet) },
}

func SortResults(rows []ProcessSummary, sortBy SortBy, order SortOrder) {
	fn, ok := sorters[sortBy]
	if !ok {
		return
	}

	slices.SortFunc(rows, func(a ProcessSummary, b ProcessSummary) int {
		result := fn(a, b)

		if order == SortOrderDescending {
			return -result
		}

		return result
	})
}
