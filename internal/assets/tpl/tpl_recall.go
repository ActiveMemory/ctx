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

	// RecallErrorMarker is the error indicator for tool results.
	RecallErrorMarker = "❌ Error"

	// RecallDetailsSummary formats the summary text for collapsible content.
	// Args: line count.
	RecallDetailsSummary = "%d lines"

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

	// PlanSummary is the <summary> label for a collapsible plan
	// section, rendered via the [Details] template.
	PlanSummary = "📋 Plan"

	// RecallApiError is a collapsed API error message.
	RecallApiError = "> ⚠ API error response (message omitted)"

	// RecallToolError formats a CC-level tool error.
	// Args: error message.
	RecallToolError = "> ⚠ Tool error: %s"

	// RecallSystemPrefix prefixes system-injected messages.
	RecallSystemPrefix = "[system] "
)
