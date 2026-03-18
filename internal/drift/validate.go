//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// templateFile checks if file content appears to be a template.
//
// Looks for common template markers like YOUR_, {{, REPLACE_, TODO, CHANGEME.
// Used to avoid flagging template files as containing secrets.
//
// Parameters:
//   - content: File content to check
//
// Returns:
//   - bool: True if content contains template markers
func templateFile(content []byte) bool {
	s := string(content)
	templateMarkers := token.TemplateMarkers
	for _, marker := range templateMarkers {
		if strings.Contains(strings.ToUpper(s), marker) {
			return true
		}
	}
	return false
}
