//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package changes

import (
	"fmt"
	"strings"
)

// RenderChanges renders the full CLI output for `ctx changes`.
func RenderChanges(refLabel string, ctxChanges []ContextChange, code CodeSummary) string {
	var b strings.Builder

	b.WriteString("## Changes Since Last Session\n\n")
	b.WriteString(fmt.Sprintf("**Reference point**: %s\n\n", refLabel))

	if len(ctxChanges) > 0 {
		b.WriteString("### Context File Changes\n")
		for _, c := range ctxChanges {
			b.WriteString(fmt.Sprintf("- `%s` — modified %s\n",
				c.Name, c.ModTime.Format("2006-01-02 15:04")))
		}
		b.WriteString("\n")
	}

	if code.CommitCount > 0 {
		b.WriteString("### Code Changes\n")
		b.WriteString(fmt.Sprintf("- **%s** since reference point\n",
			pluralize(code.CommitCount, "commit")))
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
		b.WriteString("\n")
	}

	if len(ctxChanges) == 0 && code.CommitCount == 0 {
		b.WriteString("No changes detected since reference point.\n")
	}

	return b.String()
}

// RenderChangesForHook renders a compact summary for hook injection.
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
		msg := fmt.Sprintf("%s since last session", pluralize(code.CommitCount, "commit"))
		if code.LatestMsg != "" {
			msg += fmt.Sprintf(" (latest: %s)", code.LatestMsg)
		}
		parts = append(parts, msg)
	}

	if len(parts) == 0 {
		return ""
	}

	return "Changes since last session: " + strings.Join(parts, ". ") + "\n"
}
