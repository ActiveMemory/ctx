//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	"strings"

	errPrompt "github.com/ActiveMemory/ctx/internal/err/prompt"
)

// markedSection returns the ctx block of a template: the opening
// marker through the closing marker, inclusive.
//
// Parameters:
//   - filename: Target file name, for the error message
//   - template: Full template content
//   - markerStart: Opening marker string
//   - markerEnd: Closing marker string
//
// Returns:
//   - string: The block from markerStart through markerEnd
//   - error: Non-nil if the template lacks either marker
func markedSection(
	filename, template, markerStart, markerEnd string,
) (string, error) {
	start := strings.Index(template, markerStart)
	end := strings.Index(template, markerEnd)
	if start == -1 || end == -1 {
		return "", errPrompt.TemplateMissingMarkers(filename)
	}
	return template[start : end+len(markerEnd)], nil
}
