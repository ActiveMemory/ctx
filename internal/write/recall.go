//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
)

// SkipFile prints that a file was skipped during export.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - filename: the skipped file name.
//   - reason: why it was skipped (e.g. "locked", "exists").
func SkipFile(cmd *cobra.Command, filename, reason string) {
	if cmd == nil {
		return
	}
	sprintf(cmd, "  skip %s (%s)", filename, reason)
}

// ExportedFile prints that a file was exported or updated.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - filename: the exported file name.
//   - suffix: optional annotation (e.g. "updated, frontmatter preserved").
//     Empty string omits the parenthetical.
func ExportedFile(cmd *cobra.Command, filename, suffix string) {
	if cmd == nil {
		return
	}
	if suffix != "" {
		sprintf(cmd, "  ok %s (%s)", filename, suffix)
	} else {
		sprintf(cmd, "  ok %s", filename)
	}
}

// NoSessionsForProject prints guidance when no sessions are found.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - allProjects: if true, show the generic message; otherwise suggest --all-projects.
func NoSessionsForProject(cmd *cobra.Command, allProjects bool) {
	if cmd == nil {
		return
	}
	if allProjects {
		cmd.Println("No sessions found.")
	} else {
		cmd.Println("No sessions found for this project. Use --all-projects to see all.")
	}
}

// NoSessionsWithHint prints that no sessions were found with storage hint.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - allProjects: if true, show storage path; otherwise suggest --all-projects.
func NoSessionsWithHint(cmd *cobra.Command, allProjects bool) {
	if cmd == nil {
		return
	}
	if allProjects {
		cmd.Println("No sessions found.")
		cmd.Println("")
		cmd.Println("Sessions are stored in ~/.claude/projects/")
	} else {
		cmd.Println("No sessions found for this project.")
		cmd.Println("Use --all-projects to see sessions from all projects.")
	}
}

// AmbiguousSessionMatch prints a list of matching sessions to stderr.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - query: the ambiguous query string.
//   - lines: pre-formatted lines describing each match.
func AmbiguousSessionMatch(cmd *cobra.Command, query string, lines []string) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(fmt.Sprintf("Multiple sessions match '%s':", query))
	for _, line := range lines {
		cmd.PrintErrln(line)
	}
}

// AmbiguousSessionMatchWithHint prints matching sessions with a specific-ID hint.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - query: the ambiguous query string.
//   - lines: pre-formatted lines describing each match.
//   - hint: suggested more-specific ID.
func AmbiguousSessionMatchWithHint(cmd *cobra.Command, query string, lines []string, hint string) {
	if cmd == nil {
		return
	}
	_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Multiple sessions match '%s':\n", query)
	for _, line := range lines {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "  %s\n", line)
	}
	_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "\nUse a more specific ID (e.g., ctx recall show %s)\n", hint)
}

// Aborted prints that an operation was aborted.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func Aborted(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println("Aborted.")
}

// ExportFinalSummary prints the final export summary with counts.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - exported: number of new files written.
//   - updated: number of existing files updated.
//   - renamed: number of files renamed.
//   - skipped: number of files skipped.
func ExportFinalSummary(cmd *cobra.Command, exported, updated, renamed, skipped int) {
	if cmd == nil {
		return
	}
	cmd.Println()
	if exported > 0 {
		sprintf(cmd, "Exported %d new session(s)", exported)
	}
	if updated > 0 {
		sprintf(cmd, "Updated %d existing session(s) (YAML frontmatter preserved)", updated)
	}
	if renamed > 0 {
		sprintf(cmd, "Renamed %d session(s) to title-based filenames", renamed)
	}
	if skipped > 0 {
		sprintf(cmd, "Skipped %d existing file(s).", skipped)
	}
}

// NoFiltersMatch prints that no sessions matched the applied filters.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func NoFiltersMatch(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println("No sessions match the filters.")
}

// SessionListHeader prints the session count header for recall list.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - total: total sessions found.
//   - shown: filtered count (0 to omit the parenthetical).
func SessionListHeader(cmd *cobra.Command, total, shown int) {
	if cmd == nil {
		return
	}
	if shown > 0 && shown != total {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Found %d sessions (%d shown)\n\n", total, shown)
	} else {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Found %d sessions\n\n", total)
	}
}

// SessionListRow prints a formatted row in the session list table.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - format: printf format string for the row.
//   - values: column values.
func SessionListRow(cmd *cobra.Command, format string, values ...any) {
	if cmd == nil {
		return
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), format, values...)
}

// SessionListFooter prints the footer hint for recall list.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - hasMore: if true, show the --limit hint.
func SessionListFooter(cmd *cobra.Command, hasMore bool) {
	if cmd == nil {
		return
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout())
	if hasMore {
		cmd.Println("Use --limit to see more sessions")
	}
}

// SessionDetail prints a labeled metadata line to stdout.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - label: bold metadata prefix (e.g. "**ID**:").
//   - value: the value to display.
func SessionDetail(cmd *cobra.Command, label, value string) {
	if cmd == nil {
		return
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", label, value)
}

// SessionDetailInt prints a labeled integer metadata line to stdout.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - label: bold metadata prefix.
//   - value: the integer value.
func SessionDetailInt(cmd *cobra.Command, label string, value int) {
	if cmd == nil {
		return
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s %d\n", label, value)
}

// SectionHeader prints a Markdown section heading to stdout.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - level: heading level (e.g. 1 for "#", 2 for "##").
//   - title: the heading text.
func SectionHeader(cmd *cobra.Command, level int, title string) {
	if cmd == nil {
		return
	}
	prefix := strings.Repeat("#", level)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", prefix, title)
	_, _ = fmt.Fprintln(cmd.OutOrStdout())
}

// BlankLine prints an empty line to stdout.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func BlankLine(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout())
}

// ConversationTurn prints a conversation turn header.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - index: 1-based turn number.
//   - role: display role label (e.g. "User", "Assistant").
//   - timestamp: formatted time string.
func ConversationTurn(cmd *cobra.Command, index int, role, timestamp string) {
	if cmd == nil {
		return
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "### %d. [%s] (%s)\n", index, role, timestamp)
	_, _ = fmt.Fprintln(cmd.OutOrStdout())
}

// TextBlock prints a text block followed by a blank line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - text: the text content to print.
func TextBlock(cmd *cobra.Command, text string) {
	if cmd == nil {
		return
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), text)
	_, _ = fmt.Fprintln(cmd.OutOrStdout())
}

// CodeBlock prints content wrapped in a fenced code block.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - content: the code content.
func CodeBlock(cmd *cobra.Command, content string) {
	if cmd == nil {
		return
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "```\n%s\n```\n", content)
}

// ListItem prints a Markdown list item to stdout.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - format: printf format string for the item text.
//   - args: format arguments.
func ListItem(cmd *cobra.Command, format string, args ...any) {
	if cmd == nil {
		return
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "- "+format+config.NewlineLF, args...)
}

// NumberedItem prints a numbered item to stdout.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - n: the item number.
//   - text: the item text.
func NumberedItem(cmd *cobra.Command, n int, text string) {
	if cmd == nil {
		return
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%d. %s\n", n, text)
}

// MoreTurns prints the "and N more turns" continuation line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - remaining: number of remaining turns.
func MoreTurns(cmd *cobra.Command, remaining int) {
	if cmd == nil {
		return
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "... and %d more turns\n", remaining)
}

// Hint prints a usage hint to stdout.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - text: the hint text.
func Hint(cmd *cobra.Command, text string) {
	if cmd == nil {
		return
	}
	cmd.Println(text)
}
