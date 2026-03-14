//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// RenderChanges renders the full CLI output for `ctx changes`.
//
// Parameters:
//   - refLabel: Human-readable reference time label
//   - ctxChanges: Context file changes since reference time
//   - code: Code change summary since reference time
//
// Returns:
//   - string: Formatted Markdown output
func RenderChanges(
	refLabel string, ctxChanges []ContextChange, code CodeSummary,
) string {
	var b strings.Builder

	b.WriteString("## Changes Since Last Session\n\n")
	b.WriteString(fmt.Sprintf("**Reference point**: %s\n\n", refLabel))

	if len(ctxChanges) > 0 {
		b.WriteString("### Context File Changes\n")
		for _, c := range ctxChanges {
			b.WriteString(fmt.Sprintf("- `%s` — modified %s\n",
				c.Name, c.ModTime.Format(time.DateTimeFormat)))
		}
		b.WriteString(token.NewlineLF)
	}

	if code.CommitCount > 0 {
		b.WriteString("### Code Changes\n")
		b.WriteString(fmt.Sprintf("- **%s** since reference point\n",
			Pluralize(code.CommitCount, "commit")))
		if code.LatestMsg != "" {
			b.WriteString(fmt.Sprintf("- **Latest**: %s\n", code.LatestMsg))
		}
		if len(code.Dirs) > 0 {
			b.WriteString(fmt.Sprintf("- **Directories touched**: %s\n",
				strings.Join(code.Dirs, ", ")))
		}
		if len(code.Authors) > 0 {
			b.WriteString(fmt.Sprintf("- **Authors**: %s\n",
				strings.Join(code.Authors, ", ")))
		}
		b.WriteString(token.NewlineLF)
	}

	if len(ctxChanges) == 0 && code.CommitCount == 0 {
		b.WriteString("No changes detected since reference point.\n")
	}

	return b.String()
}

// RenderChangesForHook renders a compact summary for hook injection.
//
// Parameters:
//   - refLabel: Human-readable reference time label
//   - ctxChanges: Context file changes since reference time
//   - code: Code change summary since reference time
//
// Returns:
//   - string: Compact single-line summary, or empty if no changes
func RenderChangesForHook(refLabel string, ctxChanges []ContextChange, code CodeSummary) string {
	var parts []string

	if len(ctxChanges) > 0 {
		names := make([]string, len(ctxChanges))
		for i, c := range ctxChanges {
			names[i] = c.Name
		}
		parts = append(parts, fmt.Sprintf("Context files changed (%s): %s",
			refLabel, strings.Join(names, ", ")))
	}

	if code.CommitCount > 0 {
		msg := fmt.Sprintf("%s since last session", Pluralize(code.CommitCount, "commit"))
		if code.LatestMsg != "" {
			msg += fmt.Sprintf(" (latest: %s)", code.LatestMsg)
		}
		parts = append(parts, msg)
	}

	if len(parts) == 0 {
		return ""
	}

	return "Changes since last session: " + strings.Join(parts, ". ") + token.NewlineLF
}
