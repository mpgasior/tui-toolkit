package ui

import (
	"fmt"
	"time"
)

func formatPercentage(v float64) string {
	return fmt.Sprintf("%5.2f%%", v)
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "N/A"
	}

	if d.Hours() >= 24 {
		days := int(d.Hours() / 24)
		hours := int(d.Hours()) % 24
		return fmt.Sprintf("%dd %dh", days, hours)
	}

	return d.Round(time.Second).String()
}

func formatWorkingSet(workingSet uint64) string {
	b := float64(workingSet)
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", workingSet)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	suffix := []string{"KB", "MB", "GB", "TB"}[exp]
	return fmt.Sprintf("%.2f %s", b/float64(div), suffix)
}
