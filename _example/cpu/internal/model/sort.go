package model

import (
	"cmp"
	"slices"
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

func orderedCmp[T cmp.Ordered](a, b T, order SortOrder) int {
	r := cmp.Compare(a, b)
	if r != 0 && order == SortOrderDescending {
		return -r
	}

	return r
}

var sorters = map[SortBy]func(a, b Process, order SortOrder) int{
	SortByPID: func(a, b Process, order SortOrder) int {
		return orderedCmp(a.PID, b.PID, order)
	},
	SortByName: func(a, b Process, order SortOrder) int {
		return orderedCmp(a.Name, b.Name, order)
	},
	SortByAge: func(a, b Process, order SortOrder) int {
		if a.AgeReady != b.AgeReady {
			if a.AgeReady {
				return -1
			}
			return 1
		}
		return orderedCmp(a.Age, b.Age, order)
	},
	SortByAvgCPU: func(a, b Process, order SortOrder) int {
		if a.CPUReady != b.CPUReady {
			if a.CPUReady {
				return -1
			}

			return 1
		}

		return orderedCmp(a.CPUAvg, b.CPUAvg, order)
	},
	SortByCPU: func(a, b Process, order SortOrder) int {
		if a.CPUReady != b.CPUReady {
			if a.CPUReady {
				return -1
			}

			return 1
		}

		return orderedCmp(a.CPU, b.CPU, order)
	},
	SortByMem: func(a, b Process, order SortOrder) int {
		if a.MemReady != b.MemReady {
			if a.MemReady {
				return -1
			}
			return 1
		}
		return orderedCmp(a.MemoryRSS, b.MemoryRSS, order)
	},
	SortByMaxMem: func(a, b Process, order SortOrder) int {
		if a.MemReady != b.MemReady {
			if a.MemReady {
				return -1
			}
			return 1
		}
		return orderedCmp(a.MemoryMax, b.MemoryMax, order)
	},
}

func SortResults(rows []Process, sortBy SortBy, order SortOrder) {
	fn, ok := sorters[sortBy]
	if !ok {
		return
	}

	slices.SortFunc(rows, func(a, b Process) int {
		return sorters[SortByPID](a, b, order)
	})

	slices.SortStableFunc(rows, func(a, b Process) int {
		result := fn(a, b, order)
		return result
	})
}
