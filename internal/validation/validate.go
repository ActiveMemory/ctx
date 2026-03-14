//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validation

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/token"
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
	s = regex.FileNameChar.ReplaceAllString(s, "-")
	// Remove leading/trailing hyphens
	s = strings.Trim(s, token.Dash)
	// Convert to lowercase
	s = strings.ToLower(s)
	// Limit length
	if len(s) > file.MaxNameLen {
		s = s[:file.MaxNameLen]
	}
	if s == "" {
		s = session.DefaultSessionFilename
	}
	return s
}
