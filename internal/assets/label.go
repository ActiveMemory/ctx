//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

// Bold metadata field prefixes in journal/session Markdown.
const (
	// MetadataID is the bold ID field prefix.
	MetadataID = "**ID**:"
	// MetadataDate is the bold date field prefix.
	MetadataDate = "**Date**:"
	// MetadataTime is the bold time field prefix.
	MetadataTime = "**Time**:"
	// MetadataDuration is the bold duration field prefix.
	MetadataDuration = "**Duration**:"
	// MetadataTool is the bold tool field prefix.
	MetadataTool = "**Tool**:"
	// MetadataProject is the bold project field prefix.
	MetadataProject = "**Project**:"
	// MetadataBranch is the bold branch field prefix.
	MetadataBranch = "**Branch**:"
	// MetadataModel is the bold model field prefix.
	MetadataModel = "**Model**:"
	// MetadataTurns is the bold turns field prefix.
	MetadataTurns = "**Turns**:"
	// MetadataParts is the bold parts field prefix.
	MetadataParts = "**Parts**:"
	// MetadataType is the bold type field prefix.
	MetadataType = "**Type**:"
	// MetadataStartTime is the bold start_time field prefix.
	MetadataStartTime = "**start_time**:"
	// MetadataEndTime is the bold end_time field prefix.
	MetadataEndTime = "**end_time**:"
	// MetadataSource is the bold source field prefix.
	MetadataSource = "**Source**:"
)

// Conversation role display labels used in exported journal entries.
const (
	// RoleUser is the display label for user turns.
	RoleUser = "User"
	// LabelRoleAssistant is the display label for assistant turns.
	LabelRoleAssistant = "Assistant"
)

// Journal content markers for detecting session modes.
const (
	// LabelSuggestionMode identifies suggestion mode sessions in journal content.
	LabelSuggestionMode = "SUGGESTION MODE:"
)

// YAML frontmatter field names used in journal entries.
const (
	// FrontmatterTitle is the YAML frontmatter key for the entry title.
	FrontmatterTitle = "title"
	// FrontmatterDate is the YAML frontmatter key for the entry date.
	FrontmatterDate = "date"
	// FrontmatterType is the YAML frontmatter key for the session type.
	FrontmatterType = "type"
	// FrontmatterOutcome is the YAML frontmatter key for the session outcome.
	FrontmatterOutcome = "outcome"
	// FrontmatterTopics is the YAML frontmatter key for the topics list.
	FrontmatterTopics = "topics"
	// FrontmatterTechnologies is the YAML frontmatter key for the technologies list.
	FrontmatterTechnologies = "technologies"
	// FrontmatterKeyFiles is the YAML frontmatter key for the key files list.
	FrontmatterKeyFiles = "key_files"
	// FrontmatterLocked is the YAML frontmatter key and journal state
	// marker for locked entries.
	FrontmatterLocked = "locked"
	// Unlocked is the display label for unlocked entries.
	Unlocked = "unlocked"
)

// Additional bold metadata field prefixes for session show output.
const (
	// MetadataStarted is the bold started field prefix.
	MetadataStarted = "**Started**:"
	// MetadataMessages is the bold messages field prefix.
	MetadataMessages = "**Messages**:"
	// MetadataInputUsage is the bold input token usage field prefix.
	MetadataInputUsage = "**Tokens In**:"
	// MetadataOutputUsage is the bold output token usage field prefix.
	MetadataOutputUsage = "**Tokens Out**:"
	// MetadataTotal is the bold total field prefix.
	MetadataTotal = "**Total**:"
)

// Column header labels for recall list output.
const (
	// ColSlug is the column header for session slugs.
	ColSlug = "Slug"
	// ColProject is the column header for project names.
	ColProject = "Project"
	// ColDate is the column header for dates.
	ColDate = "Date"
	// ColDuration is the column header for durations.
	ColDuration = "Duration"
	// ColTurns is the column header for turn counts.
	ColTurns = "Turns"
	// ColTokens is the column header for token counts.
	ColTokens = "Tokens"
)

// Claude Code tool names used in session transcripts.
const (
	ToolRead      = "Read"
	ToolWrite     = "Write"
	ToolEdit      = "Edit"
	ToolBash      = "Bash"
	ToolGrep      = "Grep"
	ToolGlob      = "Glob"
	ToolWebFetch  = "WebFetch"
	ToolWebSearch = "WebSearch"
	ToolTask      = "Task"
)

// Plain-text metadata labels used in HTML table rows.
const (
	MetaLabelID       = "ID"
	MetaLabelDate     = "Date"
	MetaLabelTime     = "Time"
	MetaLabelDuration = "Duration"
	MetaLabelTool     = "Tool"
	MetaLabelProject  = "Project"
	MetaLabelBranch   = "Branch"
	MetaLabelModel    = "Model"
	MetaLabelTurns    = "Turns"
	MetaLabelTokens   = "Tokens"
	MetaLabelParts    = "Parts"
)

// YAML frontmatter field keys for journal export.
const (
	FmKeyDate      = "date"
	FmKeyTime      = "time"
	FmKeyProject   = "project"
	FmKeyBranch    = "branch"
	FmKeyModel     = "model"
	FmKeyTokensIn  = "tokens_in"
	FmKeyTokensOut = "tokens_out"
	FmKeySessionID = "session_id"
	FmKeyTitle     = "title"
)

// Claude Code tool input JSON keys for display formatting.
const (
	ToolInputFilePath    = "file_path"
	ToolInputCommand     = "command"
	ToolInputPattern     = "pattern"
	ToolInputURL         = "url"
	ToolInputQuery       = "query"
	ToolInputDescription = "description"
)

// Tool display limits.
const (
	// ToolDisplayMaxLen is the max length for tool parameter display before truncation.
	ToolDisplayMaxLen = 100
)

// CLI flag names used in multiple commands.
const (
	// FlagSince is the --since flag name.
	FlagSince = "--since"
	// FlagUntil is the --until flag name.
	FlagUntil = "--until"
)

// Export action reasons for skip/export output.
const (
	// ReasonExists is the skip reason when a file already exists.
	ReasonExists = "exists"
	// ReasonUpdated is the annotation for updated files with preserved frontmatter.
	ReasonUpdated = "updated, frontmatter preserved"
)

// Section headers used in recall show output.
const (
	// SectionToolUsage is the heading for the tool usage summary.
	SectionToolUsage = "Tool Usage"
	// SectionConversation is the heading for the full conversation.
	SectionConversation = "Conversation"
	// SectionConversationPreview is the heading for the conversation preview.
	SectionConversationPreview = "Conversation Preview"
)

// Recall show guidance hints.
const (
	// HintUseFullFlag is the hint to use --full for all messages.
	HintUseFullFlag = "Use --full to see all messages"
	// HintUseAllProjects is the hint when no sessions found for the current project.
	HintUseAllProjects = "use --all-projects to search all"
)

// Inline labels for conversation output.
const (
	// LabelTool is the prefix for tool use lines in conversation output.
	LabelTool = "Tool:"
	// LabelError is the prefix for error lines in conversation output.
	LabelError = "Error:"
)

// Journal turn markers for content transformation.
const (
	// BoldReminder is the bold-style system reminder prefix.
	BoldReminder = "**System Reminder**:"
	// ToolOutput is the turn role label for tool output turns.
	ToolOutput = "Tool Output"
)

// Loop output markers.
const (
	// LoopComplete is the banner printed when the loop finishes.
	LoopComplete = "=== Loop Complete ==="
)
