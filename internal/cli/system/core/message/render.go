//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package message

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/ActiveMemory/ctx/internal/config/session"
)

// renderTemplate executes a Go text/template with the given vars.
// Returns the fallback on any parse or execution error. Returns empty
// string if the template content is empty or whitespace-only
// (intentional silence).
//
// Parameters:
//   - tmpl: Go text/template source string
//   - vars: Key-value pairs available inside the template
//   - fallback: Returned when parsing or execution fails
//
// Returns:
//   - string: Rendered output, empty string for silent
//     templates, or fallback on error
func renderTemplate(tmpl string, vars map[string]any, fallback string) string {
	if strings.TrimSpace(tmpl) == "" {
		return "" // intentional silence
	}

	t, parseErr := template.New(session.TemplateName).Parse(tmpl)
	if parseErr != nil {
		return fallback
	}

	var buf bytes.Buffer
	if execErr := t.Execute(&buf, vars); execErr != nil {
		return fallback
	}
	return buf.String()
}
