//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// loadMessage loads a hook message template by hook name and variant.
//
// Priority:
//  1. .context/hooks/messages/{hook}/{variant}.txt (user override)
//  2. internal/assets/hooks/messages/{hook}/{variant}.txt (embedded default)
//  3. fallback string (hardcoded, belt and suspenders)
//
// Returns empty string if the resolved template is empty or whitespace-only
// (intentional silence). The vars map provides template variables;
// nil is valid when no dynamic content is needed.
func loadMessage(hook, variant string, vars map[string]any, fallback string) string {
	filename := variant + ".txt"

	// 1. User override in .context/
	userPath := filepath.Join(rc.ContextDir(), "hooks", "messages", hook, filename)
	if data, readErr := os.ReadFile(userPath); readErr == nil { //nolint:gosec // project-local override path
		return renderTemplate(string(data), vars, fallback)
	}

	// 2. Embedded default
	if data, readErr := assets.HookMessage(hook, filename); readErr == nil {
		return renderTemplate(string(data), vars, fallback)
	}

	// 3. Hardcoded fallback
	return renderTemplate(fallback, vars, fallback)
}

// renderTemplate executes a Go text/template with the given vars.
// Returns the fallback on any parse or execution error. Returns empty
// string if the template content is empty or whitespace-only
// (intentional silence).
func renderTemplate(tmpl string, vars map[string]any, fallback string) string {
	if strings.TrimSpace(tmpl) == "" {
		return "" // intentional silence
	}

	t, parseErr := template.New("msg").Parse(tmpl)
	if parseErr != nil {
		return fallback
	}

	var buf bytes.Buffer
	if execErr := t.Execute(&buf, vars); execErr != nil {
		return fallback
	}
	return buf.String()
}

// boxBottom is the standard bottom border for hook message boxes.
const boxBottom = "└──────────────────────────────────────────────────"

// variantBoth is the template variant name used when both ceremony
// conditions are unmet (e.g. neither remember nor wrapup done).
const variantBoth = "both"

// sessionUnknown is the fallback session ID used when input lacks one.
const sessionUnknown = "unknown"

// boxLines wraps each line of content with the │ box-drawing prefix.
// Trailing newlines on content are trimmed before splitting to avoid
// an empty trailing box line.
func boxLines(content string) string {
	var b strings.Builder
	for _, line := range strings.Split(strings.TrimRight(content, "\n"), "\n") {
		b.WriteString("│ ")
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}
