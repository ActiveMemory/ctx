package validation

import (
	"regexp"
	"strings"
)

// SanitizeFilename converts a topic string to a safe filename component.
func SanitizeFilename(s string) string {
	// Replace spaces and special chars with hyphens
	re := regexp.MustCompile(`[^a-zA-Z0-9-]+`)
	s = re.ReplaceAllString(s, "-")
	// Remove leading/trailing hyphens
	s = strings.Trim(s, "-")
	// Convert to lowercase
	s = strings.ToLower(s)
	// Limit length
	if len(s) > 50 {
		s = s[:50]
	}
	if s == "" {
		s = "session"
	}
	return s
}
