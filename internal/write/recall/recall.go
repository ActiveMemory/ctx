//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/spf13/cobra"
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
	cmd.Println(fmt.Sprintf("  skip %s (%s)", filename, reason))
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
		cmd.Println(fmt.Sprintf("  ok %s (%s)", filename, suffix))
	} else {
		cmd.Println(fmt.Sprintf("  ok %s", filename))
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
		cmd.Println()
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
		cmd.Println(fmt.Sprintf("Exported %d new session(s)", exported))
	}
	if updated > 0 {
		cmd.Println(fmt.Sprintf("Updated %d existing session(s) (YAML frontmatter preserved)", updated))
	}
	if renamed > 0 {
		cmd.Println(fmt.Sprintf("Renamed %d session(s) to title-based filenames", renamed))
	}
	if skipped > 0 {
		cmd.Println(fmt.Sprintf("Skipped %d existing file(s).", skipped))
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

// SessionInfo holds pre-formatted session metadata for display.
type SessionInfo struct {
	Slug      string
	ID        string
	Tool      string
	Project   string
	Branch    string // empty to omit
	Model     string // empty to omit
	Started   string
	Duration  string
	Turns     int
	Messages  int
	TokensIn  string
	TokensOut string
	TokensAll string
}

// SessionMetadata prints the full session metadata block: identity,
// timing, and token usage sections separated by blank lines.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - info: pre-formatted session metadata.
func SessionMetadata(cmd *cobra.Command, info SessionInfo) {
	if cmd == nil {
		return
	}
	SectionHeader(cmd, 1, info.Slug)

	SessionDetail(cmd, assets.MetadataID, info.ID)
	SessionDetail(cmd, assets.MetadataTool, info.Tool)
	SessionDetail(cmd, assets.MetadataProject, info.Project)
	if info.Branch != "" {
		SessionDetail(cmd, assets.MetadataBranch, info.Branch)
	}
	if info.Model != "" {
		SessionDetail(cmd, assets.MetadataModel, info.Model)
	}
	BlankLine(cmd)

	SessionDetail(cmd, assets.MetadataStarted, info.Started)
	SessionDetail(cmd, assets.MetadataDuration, info.Duration)
	SessionDetailInt(cmd, assets.MetadataTurns, info.Turns)
	SessionDetailInt(cmd, assets.MetadataMessages, info.Messages)
	BlankLine(cmd)

	SessionDetail(cmd, assets.MetadataInputUsage, info.TokensIn)
	SessionDetail(cmd, assets.MetadataOutputUsage, info.TokensOut)
	SessionDetail(cmd, assets.MetadataTotal, info.TokensAll)
	BlankLine(cmd)
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
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "- "+format+token.NewlineLF, args...)
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

// LockUnlockNone prints the message when no journal entries are found (lock/unlock context).
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func LockUnlockNone(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(assets.TextDesc(assets.TextDescKeyWriteJournalSyncNone))
}

// LockUnlockEntry prints the confirmation for a single locked/unlocked entry.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - filename: journal filename.
//   - verb: "locked" or "unlocked".
func LockUnlockEntry(cmd *cobra.Command, filename, verb string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteLockUnlockEntry), filename, verb))
}

// LockUnlockSummary prints the lock/unlock summary.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - verb: "locked" or "unlocked".
//   - count: number of entries changed. Zero prints no-changes message.
func LockUnlockSummary(cmd *cobra.Command, verb string, count int) {
	if cmd == nil {
		return
	}
	if count == 0 {
		cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteLockUnlockNoChanges), verb))
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteLockUnlockSummary), strings.Title(verb), count)) //nolint:staticcheck // strings.Title is fine for single words
}

// JournalSyncNone prints the message when no journal entries are found.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func JournalSyncNone(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(assets.TextDesc(assets.TextDescKeyWriteJournalSyncNone))
}

// JournalSyncLocked prints a single locked-entry confirmation.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - filename: the journal filename that was locked.
func JournalSyncLocked(cmd *cobra.Command, filename string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteJournalSyncLocked), filename))
}

// JournalSyncUnlocked prints a single unlocked-entry confirmation.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - filename: the journal filename that was unlocked.
func JournalSyncUnlocked(cmd *cobra.Command, filename string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteJournalSyncUnlocked), filename))
}

// JournalSyncSummary prints the sync summary: match, locked count,
// and/or unlocked count.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - locked: number of newly locked entries.
//   - unlocked: number of newly unlocked entries.
func JournalSyncSummary(cmd *cobra.Command, locked, unlocked int) {
	if cmd == nil {
		return
	}
	if locked == 0 && unlocked == 0 {
		cmd.Println(assets.TextDesc(assets.TextDescKeyWriteJournalSyncMatch))
		return
	}
	if locked > 0 {
		cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteJournalSyncLockedCount), locked))
	}
	if unlocked > 0 {
		cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteJournalSyncUnlockedCount), unlocked))
	}
}
