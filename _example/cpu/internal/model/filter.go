package model

import "strings"

func Filter(rows []Process, term string, exclude Exclude) []Process {
	term = strings.ToLower(term)

	var filtered []Process
	for _, row := range rows {
		if !row.ExitTime.IsZero() && exclude == ExcludeExited {
			continue
		}

		if row.ExitTime.IsZero() && exclude == ExcludeActive {
			continue
		}

		name := strings.ToLower(row.Name)
		if term == "" || strings.Contains(name, term) {
			filtered = append(filtered, row)
		}
	}

	return filtered
}
