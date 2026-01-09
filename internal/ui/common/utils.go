package common

import (
	"fmt"
	"time"
)

// FormatRelativeTime returns a human-readable relative time string.
func FormatRelativeTime(t time.Time) string {
	if t.IsZero() {
		return "never"
	}

	elapsed := time.Since(t)

	switch {
	case elapsed < time.Minute:
		return "just now"
	case elapsed < time.Hour:
		mins := int(elapsed.Minutes())
		return fmt.Sprintf("%dm ago", mins)
	case elapsed < 24*time.Hour:
		hours := int(elapsed.Hours())
		return fmt.Sprintf("%dh ago", hours)
	case elapsed < 7*24*time.Hour:
		days := int(elapsed.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	case elapsed < 30*24*time.Hour:
		weeks := int(elapsed.Hours() / (24 * 7))
		return fmt.Sprintf("%dw ago", weeks)
	case elapsed < 365*24*time.Hour:
		months := int(elapsed.Hours() / (24 * 30))
		return fmt.Sprintf("%dmo ago", months)
	default:
		years := int(elapsed.Hours() / (24 * 365))
		return fmt.Sprintf("%dy ago", years)
	}
}
