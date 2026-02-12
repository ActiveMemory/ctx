//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

// Recall export format templates.
//
// These templates define the structure of exported session transcripts.
// Each uses fmt.Sprintf verbs for interpolation.
const (
	// TplRecallFilename formats a journal entry filename.
	// Args: date, slug, shortID.
	TplRecallFilename = "%s-%s-%s.md"

	// TplRecallTokens formats the token stats line.
	// Args: total, in, out.
	TplRecallTokens = "**Tokens**: %s (in: %s, out: %s)" //nolint:gosec // G101: display template, not a credential

	// TplRecallPartOf formats the part indicator.
	// Args: part, totalParts.
	TplRecallPartOf = "**Part %d of %d**"

	// TplRecallConversationContinued formats the continued conversation heading.
	// Args: previous part number.
	TplRecallConversationContinued = "## Conversation (continued from part %d)"

	// TplRecallTurnHeader formats a conversation turn heading.
	// Args: msgNum, role, time.
	TplRecallTurnHeader = "### %d. %s (%s)"

	// TplRecallToolUse formats a tool use line.
	// Args: formatted tool name and args.
	TplRecallToolUse = "🔧 **%s**"

	// TplRecallToolCount formats a tool usage count line.
	// Args: name, count.
	TplRecallToolCount = "- %s: %d"

	// TplRecallSummaryPlaceholder is the placeholder text in the summary section.
	TplRecallSummaryPlaceholder = "[Add your summary of this session]"

	// TplRecallErrorMarker is the error indicator for tool results.
	TplRecallErrorMarker = "❌ Error"

	// TplRecallDetailsSummary formats the summary text for collapsible content.
	// Args: line count.
	TplRecallDetailsSummary = "%d lines"

	// TplRecallDetailsOpen formats the opening HTML for collapsible content.
	// Args: summary text.
	TplRecallDetailsOpen = "<details>\n<summary>%s</summary>"

	// TplRecallDetailsClose is the closing HTML for collapsible content.
	TplRecallDetailsClose = "</details>"

	// TplRecallFencedBlock formats content inside code fences.
	// Args: fence, content, fence.
	TplRecallFencedBlock = "%s\n%s\n%s"

	// TplRecallNavPrev formats the previous part navigation link.
	// Args: filename.
	TplRecallNavPrev = "[← Previous](%s)"

	// TplRecallNavNext formats the next part navigation link.
	// Args: filename.
	TplRecallNavNext = "[Next →](%s)"

	// TplRecallPartFilename formats a multi-part filename.
	// Args: baseName, part.
	TplRecallPartFilename = "%s-p%d.md"
)
