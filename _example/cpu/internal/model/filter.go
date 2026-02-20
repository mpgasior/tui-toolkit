package model

import "strings"

func Filter(rows []ProcessSummary, term string) []ProcessSummary {
	if term == "" {
		return rows
	}

	term = strings.ToLower(term)

	var filtered []ProcessSummary
	for _, row := range rows {
		name := strings.ToLower(row.Name)
		if strings.Contains(name, term) {
			filtered = append(filtered, row)
		}
	}

	return filtered
}
