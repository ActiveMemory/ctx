//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

// Recall export format templates.
//
// These templates define the structure of exported session transcripts.
// Each uses fmt.Sprintf verbs for interpolation.
const (
	// RecallFilename formats a journal entry filename.
	// Args: date, slug, shortID.
	RecallFilename = "%s-%s-%s.md"

	// RecallTokens formats the token stats line.
	// Args: total, in, out.
	RecallTokens = "**Tokens**: %s (in: %s, out: %s)" //nolint:gosec // G101: display template, not a credential

	// RecallPartOf formats the part indicator.
	// Args: part, totalParts.
	RecallPartOf = "**Part %d of %d**"

	// RecallConversationContinued formats the continued conversation heading.
	// Args: previous part number.
	RecallConversationContinued = "## Conversation (continued from part %d)"

	// RecallTurnHeader formats a conversation turn heading.
	// Args: msgNum, role, time.
	RecallTurnHeader = "### %d. %s (%s)"

	// RecallToolUse formats a tool use line.
	// Args: formatted tool name and args.
	RecallToolUse = "🔧 **%s**"

	// RecallToolCount formats a tool usage count line.
	// Args: name, count.
	RecallToolCount = "- %s: %d"

	// RecallSummaryPlaceholder is the placeholder text in the summary section.
	RecallSummaryPlaceholder = "[Add your summary of this session]"

	// RecallErrorMarker is the error indicator for tool results.
	RecallErrorMarker = "❌ Error"

	// RecallDetailsSummary formats the summary text for collapsible content.
	// Args: line count.
	RecallDetailsSummary = "%d lines"

	// RecallDetailsOpen formats the opening HTML for collapsible content.
	// Args: summary text. INVARIANT: the <summary> tag is always single-line
	// (<summary>N lines</summary>). Multi-line <summary> blocks (standalone
	// <summary> on its own line) are Claude Code context compaction artifacts
	// and are stripped by stripSystemReminders. This distinction is the basis
	// for safe disambiguation.
	RecallDetailsOpen = "<details>\n<summary>%s</summary>"

	// RecallDetailsClose is the closing HTML for collapsible content.
	RecallDetailsClose = "</details>"

	// RecallFencedBlock formats content inside code fences.
	// Args: fence, content, fence.
	RecallFencedBlock = "%s\n%s\n%s"

	// RecallNavPrev formats the previous part navigation link.
	// Args: filename.
	RecallNavPrev = "[← Previous](%s)"

	// RecallNavNext formats the next part navigation link.
	// Args: filename.
	RecallNavNext = "[Next →](%s)"

	// RecallPartFilename formats a multi-part filename.
	// Args: baseName, part.
	RecallPartFilename = "%s-p%d.md"

	// RecallListRow is the printf meta-format for recall list table rows.
	// Args: slugWidth, projectWidth. Produces a format string for 6 columns.
	RecallListRow = "  %%-%ds  %%-%ds  %%-17s  %%8s  %%5s  %%7s\n"

	// SessionMatch formats a session match line for ambiguous queries.
	// Args: slug, shortID, dateTime.
	SessionMatch = "%s (%s) - %s"

	// MetaDetailsOpen opens a collapsible details block with an HTML table.
	// Markdown tables don't render inside <details> in Zensical, so we use HTML.
	// Args: summary text.
	MetaDetailsOpen = "<details>\n<summary>%s</summary>\n<table>"

	// MetaDetailsClose closes a collapsible details block with HTML table.
	MetaDetailsClose = "</table>\n</details>"

	// MetaRow formats a single row in an HTML metadata table.
	// Args: label, value.
	MetaRow = "<tr><td><strong>%s</strong></td><td>%s</td></tr>"

	// FmQuoted formats a YAML frontmatter quoted string field.
	// Args: key, value.
	FmQuoted = "%s: %q"

	// FmString formats a YAML frontmatter bare string field.
	// Args: key, value.
	FmString = "%s: %s"

	// FmInt formats a YAML frontmatter integer field.
	// Args: key, value.
	FmInt = "%s: %d"

	// ToolDisplay formats a tool name with its key parameter.
	// Args: tool name, parameter value.
	ToolDisplay = "%s: %s"
)
