//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for MCP tool output.
const (
	// DescKeyMCPToolStatusDesc is the text key for mcp tool status desc messages.
	DescKeyMCPToolStatusDesc = "mcp.tool-status-desc"
	// DescKeyMCPToolAddDesc is the text key for mcp tool add desc messages.
	DescKeyMCPToolAddDesc = "mcp.tool-add-desc"
	// DescKeyMCPToolCompleteDesc is the text key for mcp tool complete desc
	// messages.
	DescKeyMCPToolCompleteDesc = "mcp.tool-complete-desc"
	// DescKeyMCPToolDriftDesc is the text key for mcp tool drift desc messages.
	DescKeyMCPToolDriftDesc = "mcp.tool-drift-desc"
	// DescKeyMCPToolJournalSourceDesc is the text key for mcp tool journal source
	// desc messages.
	DescKeyMCPToolJournalSourceDesc = "mcp.tool-journal-source-desc"
	// DescKeyMCPToolWatchUpdateDesc is the text key for mcp tool watch update
	// desc messages.
	DescKeyMCPToolWatchUpdateDesc = "mcp.tool-watch-update-desc"
	// DescKeyMCPToolCompactDesc is the text key for mcp tool compact desc
	// messages.
	DescKeyMCPToolCompactDesc = "mcp.tool-compact-desc"
	// DescKeyMCPToolNextDesc is the text key for mcp tool next desc messages.
	DescKeyMCPToolNextDesc = "mcp.tool-next-desc"
	// DescKeyMCPToolCheckTaskDesc is the text key for mcp tool check task desc
	// messages.
	DescKeyMCPToolCheckTaskDesc = "mcp.tool-check-task-desc"
	// DescKeyMCPToolSessionDesc is the text key for mcp tool session desc
	// messages.
	DescKeyMCPToolSessionDesc = "mcp.tool-session-desc"
	// DescKeyMCPToolRemindDesc is the text key for mcp tool remind desc messages.
	DescKeyMCPToolRemindDesc = "mcp.tool-remind-desc"
	// DescKeyMCPToolPropType is the text key for mcp tool prop type messages.
	DescKeyMCPToolPropType = "mcp.tool-prop-type"
	// DescKeyMCPToolPropContent is the text key for mcp tool prop content
	// messages.
	DescKeyMCPToolPropContent = "mcp.tool-prop-content"
	// DescKeyMCPToolPropPriority is the text key for mcp tool prop priority
	// messages.
	DescKeyMCPToolPropPriority = "mcp.tool-prop-priority"
	// DescKeyMCPToolPropContext is the text key for mcp tool prop context
	// messages.
	DescKeyMCPToolPropContext = "mcp.tool-prop-context"
	// DescKeyMCPToolPropRationale is the text key for mcp tool prop rationale
	// messages.
	DescKeyMCPToolPropRationale = "mcp.tool-prop-rationale"
	// DescKeyMCPToolPropConseq is the text key for mcp tool prop conseq messages.
	DescKeyMCPToolPropConseq = "mcp.tool-prop-consequence"
	// DescKeyMCPToolPropLesson is the text key for mcp tool prop lesson messages.
	DescKeyMCPToolPropLesson = "mcp.tool-prop-lesson"
	// DescKeyMCPToolPropApplication is the text key for mcp tool prop application
	// messages.
	DescKeyMCPToolPropApplication = "mcp.tool-prop-application"
	// DescKeyMCPToolPropQuery is the text key for mcp tool prop query messages.
	DescKeyMCPToolPropQuery = "mcp.tool-prop-query"
	// DescKeyMCPToolPropLimit is the text key for mcp tool prop limit messages.
	DescKeyMCPToolPropLimit = "mcp.tool-prop-limit"
	// DescKeyMCPToolPropSince is the text key for mcp tool prop since messages.
	DescKeyMCPToolPropSince = "mcp.tool-prop-since"
	// DescKeyMCPToolPropEntryType is the text key for mcp tool prop entry type
	// messages.
	DescKeyMCPToolPropEntryType = "mcp.tool-prop-entry-type"
	// DescKeyMCPToolPropMainContent is the text key for mcp tool prop main
	// content messages.
	DescKeyMCPToolPropMainContent = "mcp.tool-prop-main-content"
	// DescKeyMCPToolPropCtxBg is the text key for mcp tool prop ctx bg messages.
	DescKeyMCPToolPropCtxBg = "mcp.tool-prop-ctx-background"
	// DescKeyMCPToolPropArchive is the text key for mcp tool prop archive
	// messages.
	DescKeyMCPToolPropArchive = "mcp.tool-prop-archive"
	// DescKeyMCPToolPropRecentAct is the text key for mcp tool prop recent act
	// messages.
	DescKeyMCPToolPropRecentAct = "mcp.tool-prop-recent-action"
	// DescKeyMCPToolPropEventType is the text key for mcp tool prop event type
	// messages.
	DescKeyMCPToolPropEventType = "mcp.tool-prop-event-type"
	// DescKeyMCPToolPropCaller is the text key for mcp tool prop caller messages.
	DescKeyMCPToolPropCaller = "mcp.tool-prop-caller"
	// DescKeyMCPToolSteeringGetDesc is the text key for mcp tool steering get
	// desc messages.
	DescKeyMCPToolSteeringGetDesc = "mcp.tool-steering-get-desc"
	// DescKeyMCPToolSearchDesc is the text key for mcp tool search desc messages.
	DescKeyMCPToolSearchDesc = "mcp.tool-search-desc"
	// DescKeyMCPToolSessionStartDesc is the text key for mcp tool session start
	// desc messages.
	DescKeyMCPToolSessionStartDesc = "mcp.tool-session-start-desc"
	// DescKeyMCPToolSessionEndDesc is the text key for mcp tool session end desc
	// messages.
	DescKeyMCPToolSessionEndDesc = "mcp.tool-session-end-desc"
	// DescKeyMCPToolPropPrompt is the text key for mcp tool prop prompt messages.
	DescKeyMCPToolPropPrompt = "mcp.tool-prop-prompt"
	// DescKeyMCPToolPropSearchQuery is the text key for mcp tool prop search
	// query messages.
	DescKeyMCPToolPropSearchQuery = "mcp.tool-prop-search-query"
	// DescKeyMCPToolPropSummary is the text key for mcp tool prop summary
	// messages.
	DescKeyMCPToolPropSummary = "mcp.tool-prop-summary"
)

// DescKeys for MCP handler steering/search output.
const (
	// DescKeyMCPSteeringSection is the text key for mcp steering section messages.
	DescKeyMCPSteeringSection = "mcp.steering-section"
	// DescKeyMCPSearchHitLine is the text key for mcp search hit line messages.
	DescKeyMCPSearchHitLine = "mcp.search-hit-line"
	// DescKeyMCPSearchNoMatch is the text key for mcp search no match messages.
	DescKeyMCPSearchNoMatch = "mcp.search-no-match"
)
