package model

import "strings"

func Filter(rows []QueryResult, term string) []QueryResult {
	if term == "" {
		return rows
	}

	term = strings.ToLower(term)

	var filtered []QueryResult
	for _, row := range rows {
		name := strings.ToLower(row.Name)
		if strings.Contains(name, term) {
			filtered = append(filtered, row)
		}
	}

	return filtered
}
