package validation

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// SanitizeFilename converts a topic string to a safe filename component.
//
// Replaces spaces and special characters with hyphens, converts to lowercase,
// and limits length to 50 characters. Returns "session" if the input is empty.
//
// Parameters:
//   - s: Topic string to sanitize
//
// Returns:
//   - string: Safe filename component (lowercase, hyphenated, max 50 chars)
func SanitizeFilename(s string) string {
	// Replace spaces and special chars with hyphens
	s = config.RegExNonFileNameChar.ReplaceAllString(s, "-")
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
