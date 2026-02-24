package model

import "strings"

func Filter(rows []Process, term string) []Process {
	if term == "" {
		return rows
	}

	term = strings.ToLower(term)

	var filtered []Process
	for _, row := range rows {
		name := strings.ToLower(row.Name)
		if strings.Contains(name, term) {
			filtered = append(filtered, row)
		}
	}

	return filtered
}
