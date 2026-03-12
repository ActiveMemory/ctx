//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package file

import (
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/entry"
)

// Session defaults.
const (
	// MaxNameLen is the maximum character length for sanitized filename components.
	MaxNameLen = 50
)

// Task archive/snapshot constants.
const (
	// ArchiveScopeTasks is the scope identifier for task archives.
	ArchiveScopeTasks = "tasks"
	// DefaultSnapshotName is the default name when no snapshot name is provided.
	DefaultSnapshotName = "snapshot"
	// SnapshotFilenameFormat is the filename template for task snapshots.
	// Args: name, formatted timestamp.
	SnapshotFilenameFormat = "tasks-%s-%s" + ExtMarkdown
	// SnapshotTimeFormat is the compact timestamp layout for snapshot filenames.
	SnapshotTimeFormat = "2006-01-02-1504"
)

// Stream scanner buffer sizes.
const (
	// StreamScannerInitCap is the initial capacity for the scanner buffer.
	StreamScannerInitCap = 64 * 1024
	// StreamScannerMaxSize is the maximum size for the scanner buffer.
	StreamScannerMaxSize = 1024 * 1024
)

// Runtime configuration constants.
const (
	// FileContextRC is the optional runtime configuration file.
	FileContextRC = ".ctxrc"
)

// Environment configuration.
const (
	// EnvHome is the environment variable for the user's home directory.
	EnvHome = "HOME"
	// EnvCtxDir is the environment variable for overriding the context directory.
	EnvCtxDir = "CTX_DIR"
	// EnvCtxTokenBudget is the environment variable for overriding the token budget.
	EnvCtxTokenBudget = "CTX_TOKEN_BUDGET" //nolint:gosec // G101: env var name, not a credential

	// TaskBudgetPct is the fraction of the token budget allocated to tasks.
	TaskBudgetPct = 0.40
	// ConventionBudgetPct is the fraction of the token budget allocated to conventions.
	ConventionBudgetPct = 0.20
	// DefaultAgentCooldown is the default cooldown between agent context packet emissions.
	DefaultAgentCooldown = 10 * time.Minute
	// PrefixAgentTombstone is the filename prefix for agent cooldown tombstone files.
	PrefixAgentTombstone = "ctx-agent-"

	// RecencyDaysWeek is the threshold for "recent" entries (0-7 days).
	RecencyDaysWeek = 7
	// RecencyDaysMonth is the threshold for "this month" entries (8-30 days).
	RecencyDaysMonth = 30
	// RecencyDaysQuarter is the threshold for "this quarter" entries (31-90 days).
	RecencyDaysQuarter = 90
	// RecencyScoreWeek is the recency score for entries within a week.
	RecencyScoreWeek = 1.0
	// RecencyScoreMonth is the recency score for entries within a month.
	RecencyScoreMonth = 0.7
	// RecencyScoreQuarter is the recency score for entries within a quarter.
	RecencyScoreQuarter = 0.4
	// RecencyScoreOld is the recency score for entries older than a quarter.
	RecencyScoreOld = 0.2
	// RelevanceMatchCap is the keyword match count that yields maximum relevance (1.0).
	RelevanceMatchCap = 3

	// PrefixCtxLoaded is the filename prefix for session-loaded marker files.
	PrefixCtxLoaded = "ctx-loaded-"
	// EventContextLoadGate is the event name for context load gate hook events.
	EventContextLoadGate = "context-load-gate"
	// ContextLoadSeparatorChar is the character used for header/footer separators.
	ContextLoadSeparatorChar = "="
	// ContextLoadSeparatorWidth is the width of header/footer separator lines.
	ContextLoadSeparatorWidth = 80
	// ContextLoadIndexSuffix is the suffix appended to filenames for index entries.
	ContextLoadIndexSuffix = " (idx)"
	// JSONKeyTimestamp is the JSON key for timestamp extraction in event logs.
	JSONKeyTimestamp = `"timestamp":"`

	// TplArchiveFilename is the format for dated archive filenames.
	// Args: prefix, date.
	TplArchiveFilename = "%s-%s" + ExtMarkdown
	// ArchiveDateSep is the separator between heading and date in archive headers.
	ArchiveDateSep = " - "
	// TaskCompleteReplace is the regex replacement string for marking a task done.
	TaskCompleteReplace = "$1- [x] $3"

	// Profile file names and identifiers for .ctxrc management.
	FileCtxRC     = ".ctxrc"
	FileCtxRCBase = ".ctxrc.base"
	FileCtxRCDev  = ".ctxrc.dev"
	ProfileDev    = "dev"
	ProfileBase   = "base"
	ProfileProd   = "prod" // Alias for ProfileBase
	// ProfileDetectKey is the .ctxrc key that distinguishes dev from base profile.
	ProfileDetectKey = "notify:"

	// EnvBackupSMBURL is the environment variable for the SMB share URL.
	EnvBackupSMBURL = "CTX_BACKUP_SMB_URL"
	// EnvBackupSMBSubdir is the environment variable for the SMB share subdirectory.
	EnvBackupSMBSubdir = "CTX_BACKUP_SMB_SUBDIR"
	// EnvSkipPathCheck is the environment variable that skips the PATH
	// validation during init. Set to EnvTrue in tests.
	EnvSkipPathCheck = "CTX_SKIP_PATH_CHECK"
)

// Environment toggle values.
const (
	// EnvTrue is the canonical truthy value for environment variable toggles.
	EnvTrue = "1"
)

// User confirmation input values.
const (
	// ConfirmShort is the short affirmative response for y/N prompts.
	ConfirmShort = "y"
	// ConfirmLong is the long affirmative response for y/N prompts.
	ConfirmLong = "yes"
)

// Backup configuration.
const (
	// BackupDefaultSubdir is the default subdirectory on the SMB share.
	BackupDefaultSubdir = "ctx-sessions"
	// BackupMarkerFile is the state file touched on successful project backup.
	BackupMarkerFile = "ctx-last-backup"
	// BackupScopeProject backs up only the project context.
	BackupScopeProject = "project"
	// BackupScopeGlobal backs up only global Claude data.
	BackupScopeGlobal = "global"
	// BackupScopeAll backs up both project and global.
	BackupScopeAll = "all"
	// BackupTplProjectArchive is the filename template for project archives.
	// Argument: timestamp.
	BackupTplProjectArchive = "ctx-backup-%s.tar.gz"
	// BackupTplGlobalArchive is the filename template for global archives.
	// Argument: timestamp.
	BackupTplGlobalArchive = "claude-global-backup-%s.tar.gz"
	// BackupTimestampFormat is the compact timestamp layout for backup filenames.
	BackupTimestampFormat = "20060102-150405"
	// BackupExcludeTodos is the directory name excluded from global backups.
	BackupExcludeTodos = "todos"
	// BackupMarkerDir is the XDG state directory for the backup marker.
	BackupMarkerDir = ".local/state"
	// BackupMaxAgeDays is the threshold in days before a backup is considered stale.
	BackupMaxAgeDays = 2
	// BackupThrottleID is the state file name for daily throttle of backup age checks.
	BackupThrottleID = "backup-reminded"
	// FileBashrc is the user's bash configuration file.
	FileBashrc = ".bashrc"
)

// Hook name constants — used for LoadMessage, NewTemplateRef, notify.Send,
// and eventlog.Append to avoid magic strings.
const (
	// HookBlockDangerousCommands is the hook name for blocking dangerous commands.
	HookBlockDangerousCommands = "block-dangerous-commands"
	// HookBlockNonPathCtx is the hook name for blocking non-PATH ctx invocations.
	HookBlockNonPathCtx = "block-non-path-ctx"
	// HookCheckBackupAge is the hook name for backup staleness checks.
	HookCheckBackupAge = "check-backup-age"
	// HookCheckCeremonies is the hook name for ceremony usage checks.
	HookCheckCeremonies = "check-ceremonies"
	// HookCheckContextSize is the hook name for context window size checks.
	HookCheckContextSize = "check-context-size"
	// HookCheckJournal is the hook name for journal health checks.
	HookCheckJournal = "check-journal"
	// HookCheckKnowledge is the hook name for knowledge file health checks.
	HookCheckKnowledge = "check-knowledge"
	// HookCheckMapStaleness is the hook name for architecture map staleness checks.
	HookCheckMapStaleness = "check-map-staleness"
	// HookCheckMemoryDrift is the hook name for memory drift checks.
	HookCheckMemoryDrift = "check-memory-drift"
	// MemoryDriftThrottlePrefix is the state file prefix for per-session
	// memory drift nudge tombstones.
	MemoryDriftThrottlePrefix = "memory-drift-nudged-"
	// HookCheckPersistence is the hook name for context persistence nudges.
	HookCheckPersistence = "check-persistence"
	// HookCheckReminders is the hook name for session reminder checks.
	HookCheckReminders = "check-reminders"
	// HookCheckResources is the hook name for resource usage checks.
	HookCheckResources = "check-resources"
	// HookCheckTaskCompletion is the hook name for task completion nudges.
	HookCheckTaskCompletion = "check-task-completion"
	// HookCheckVersion is the hook name for version mismatch checks.
	HookCheckVersion = "check-version"
	// HookHeartbeat is the hook name for session heartbeat events.
	HookHeartbeat = "heartbeat"
	// HookPostCommit is the hook name for post-commit nudges.
	HookPostCommit = "post-commit"
	// HookQAReminder is the hook name for QA reminder gates.
	HookQAReminder = "qa-reminder"
	// HookSpecsNudge is the hook name for specs directory nudges.
	HookSpecsNudge = "specs-nudge"
	// HookVersionDrift is the hook name for version drift nudges.
	HookVersionDrift = "version-drift"
)

// Hook event names (Claude Code hook lifecycle stages).
const (
	// HookEventPreToolUse is the hook event for pre-tool-use hooks.
	HookEventPreToolUse = "PreToolUse"
	// HookEventPostToolUse is the hook event for post-tool-use hooks.
	HookEventPostToolUse = "PostToolUse"
)

// Notification channel names.
const (
	// NotifyChannelHeartbeat is the notification channel for heartbeat events.
	NotifyChannelHeartbeat = "heartbeat"
	// NotifyChannelNudge is the notification channel for nudge messages.
	NotifyChannelNudge = "nudge"
	// NotifyChannelRelay is the notification channel for relay messages.
	NotifyChannelRelay = "relay"
)

// Bootstrap display constants.
const (
	// BootstrapFileListWidth is the character width at which the file list wraps.
	BootstrapFileListWidth = 55
	// BootstrapFileListIndent is the indentation prefix for file list lines.
	BootstrapFileListIndent = "  "
)

// Task parsing constants.
const (
	// SubTaskMinIndent is the minimum indent length (in spaces) for a line
	// to be considered a subtask rather than a top-level task.
	SubTaskMinIndent = 2
)

// Numbered list parsing constants.
const (
	// NumberedListSep is the separator between the number and text in numbered lists (e.g. "1. item").
	NumberedListSep = ". "
	// NumberedListMaxDigits is the maximum index position for the separator to be recognized as a prefix.
	NumberedListMaxDigits = 2
)

// Hook decision constants — JSON values returned by PreToolUse hooks.
const (
	// HookDecisionBlock is the decision value that prevents tool execution.
	HookDecisionBlock = "block"
)

// Hook variant constants — template selectors passed to LoadMessage and
// NewTemplateRef to choose the appropriate message for each trigger type.
const (
	// VariantMidSudo selects the mid-command sudo block message.
	VariantMidSudo = "mid-sudo"
	// VariantMidGitPush selects the mid-command git push block message.
	VariantMidGitPush = "mid-git-push"
	// VariantCpToBin selects the cp/mv to bin block message.
	VariantCpToBin = "cp-to-bin"
	// VariantInstallToLocalBin selects the install to ~/.local/bin block message.
	VariantInstallToLocalBin = "install-to-local-bin"
	// VariantDotSlash selects the relative path (./ctx) block message.
	VariantDotSlash = "dot-slash"
	// VariantGoRun selects the go run block message.
	VariantGoRun = "go-run"
	// VariantAbsolutePath selects the absolute path block message.
	VariantAbsolutePath = "absolute-path"
	// VariantBoth selects the template for both ceremonies missing.
	VariantBoth = "both"
	// VariantRemember selects the template for missing /ctx-remember.
	VariantRemember = "remember"
	// VariantWrapup selects the template for missing /ctx-wrap-up.
	VariantWrapup = "wrapup"
	// VariantUnexported selects the unexported journal entries variant.
	VariantUnexported = "unexported"
	// VariantUnenriched selects the unenriched journal entries variant.
	VariantUnenriched = "unenriched"
	// VariantWarning selects the generic warning variant.
	VariantWarning = "warning"
	// VariantAlert selects the alert variant.
	VariantAlert = "alert"
	// VariantBilling selects the billing threshold variant.
	VariantBilling = "billing"
	// VariantCheckpoint selects the checkpoint variant.
	VariantCheckpoint = "checkpoint"
	// VariantGate selects the gate variant.
	VariantGate = "gate"
	// VariantKeyRotation selects the key rotation variant.
	VariantKeyRotation = "key-rotation"
	// VariantMismatch selects the version mismatch variant.
	VariantMismatch = "mismatch"
	// VariantNudge selects the generic nudge variant.
	VariantNudge = "nudge"
	// VariantOversize selects the oversize threshold variant.
	VariantOversize = "oversize"
	// VariantPulse selects the heartbeat pulse variant.
	VariantPulse = "pulse"
	// VariantReminders selects the reminders variant.
	VariantReminders = "reminders"
	// VariantStale selects the staleness variant.
	VariantStale = "stale"
	// VariantWindow selects the context window variant.
	VariantWindow = "window"
)

// Template variable key constants — used as map keys in template.Execute
// data maps to avoid magic strings in hook and display code.
const (
	// TplVarAlertMessages is the template variable for resource alert messages.
	TplVarAlertMessages = "AlertMessages"

	// TplVarUnenrichedCount is the template variable for unenriched entry count.
	TplVarUnenrichedCount = "UnenrichedCount"

	// TplVarUnexportedCount is the template variable for unexported session count.
	TplVarUnexportedCount = "UnexportedCount"

	// TplVarBinaryVersion is the template variable for the binary version string.
	TplVarBinaryVersion = "BinaryVersion"

	// TplVarFileWarnings is the template variable for knowledge file warnings.
	TplVarFileWarnings = "FileWarnings"

	// TplVarKeyAgeDays is the template variable for API key age in days.
	TplVarKeyAgeDays = "KeyAgeDays"

	// TplVarLastRefreshDate is the template variable for the last map refresh date.
	TplVarLastRefreshDate = "LastRefreshDate"

	// TplVarModuleCount is the template variable for the number of changed modules.
	TplVarModuleCount = "ModuleCount"

	// TplVarPercentage is the template variable for context window percentage.
	TplVarPercentage = "Percentage"

	// TplVarPluginVersion is the template variable for the plugin version string.
	TplVarPluginVersion = "PluginVersion"

	// TplVarPromptCount is the template variable for the prompt counter.
	TplVarPromptCount = "PromptCount"

	// TplVarPromptsSinceNudge is the template variable for prompts since last nudge.
	TplVarPromptsSinceNudge = "PromptsSinceNudge"

	// TplVarReminderList is the template variable for formatted reminder list.
	TplVarReminderList = "ReminderList"

	// TplVarThreshold is the template variable for a token threshold value.
	TplVarThreshold = "Threshold"

	// TplVarTokenCount is the template variable for a token count value.
	TplVarTokenCount = "TokenCount"

	// TplVarWarnings is the template variable for backup warning messages.
	TplVarWarnings = "Warnings"

	// TplVarHeartbeatPromptCount is the heartbeat field for prompt count.
	TplVarHeartbeatPromptCount = "prompt_count"
	// TplVarHeartbeatSessionID is the heartbeat field for session identifier.
	TplVarHeartbeatSessionID = "session_id"
	// TplVarHeartbeatContextModified is the heartbeat field for context modification flag.
	TplVarHeartbeatContextModified = "context_modified"
	// TplVarHeartbeatTokens is the heartbeat field for token count.
	TplVarHeartbeatTokens = "tokens"
	// TplVarHeartbeatContextWindow is the heartbeat field for context window size.
	TplVarHeartbeatContextWindow = "context_window"
	// TplVarHeartbeatUsagePct is the heartbeat field for usage percentage.
	TplVarHeartbeatUsagePct = "usage_pct"
)

// Auto-prune configuration.
const (
	// HoursPerDay is the number of hours in a day for duration calculations.
	HoursPerDay = 24
	// AutoPruneStaleDays is the number of days after which session state
	// files are eligible for auto-pruning during context load.
	AutoPruneStaleDays = 7
)

// Stats display configuration.
const (
	// StatsFilePrefix is the filename prefix for per-session stats JSONL files.
	StatsFilePrefix = "stats-"
	// StatsReadBufSize is the byte buffer size for reading new lines
	// from stats files during follow/stream mode.
	StatsReadBufSize = 8192
	// StatsHeaderTime is the column header label for timestamp.
	StatsHeaderTime = "TIME"
	// StatsHeaderSession is the column header label for session ID.
	StatsHeaderSession = "SESSION"
	// StatsHeaderPrompt is the column header label for prompt count.
	StatsHeaderPrompt = "PROMPT"
	// StatsHeaderTokens is the column header label for token count.
	StatsHeaderTokens = "TOKENS"
	// StatsHeaderPct is the column header label for percentage.
	StatsHeaderPct = "PCT"
	// StatsHeaderEvent is the column header label for event type.
	StatsHeaderEvent = "EVENT"
	// StatsSepTime is the column separator for the time field.
	StatsSepTime = "-------------------"
	// StatsSepSession is the column separator for the session field.
	StatsSepSession = "--------"
	// StatsSepPrompt is the column separator for the prompt field.
	StatsSepPrompt = "------"
	// StatsSepTokens is the column separator for the tokens field.
	StatsSepTokens = "--------"
	// StatsSepPct is the column separator for the percentage field.
	StatsSepPct = "----"
	// StatsSepEvent is the column separator for the event field.
	StatsSepEvent = "------------"
)

// Events display configuration.
const (
	// EventsMessageMaxLen is the maximum character length for event messages
	// in human-readable output before truncation.
	EventsMessageMaxLen = 60
	// EventsHookFallback is the placeholder displayed when no hook name
	// can be determined from an event payload.
	EventsHookFallback = "-"
	// EventsTruncationSuffix is appended to truncated event messages.
	EventsTruncationSuffix = "..."
)

// Heartbeat state file prefixes.
const (
	// HeartbeatCounterPrefix is the state file prefix for per-session
	// heartbeat prompt counters.
	HeartbeatCounterPrefix = "heartbeat-"
	// HeartbeatMtimePrefix is the state file prefix for per-session
	// heartbeat context mtime tracking.
	HeartbeatMtimePrefix = "heartbeat-mtime-"
	// HeartbeatLogFile is the log filename for heartbeat events.
	HeartbeatLogFile = "heartbeat.log"
)

// Message table formatting.
const (
	// MessageColHook is the column width for the Hook field in message list output.
	MessageColHook = 24
	// MessageColVariant is the column width for the Variant field in message list output.
	MessageColVariant = 20
	// MessageColCategory is the column width for the Category field in message list output.
	MessageColCategory = 16
	// MessageSepHook is the separator width for the Hook column underline.
	MessageSepHook = 22
	// MessageSepVariant is the separator width for the Variant column underline.
	MessageSepVariant = 18
	// MessageSepCategory is the separator width for the Category column underline.
	MessageSepCategory = 14
	// MessageSepOverride is the separator width for the Override column underline.
	MessageSepOverride = 8
)

// Resources display formatting.
const (
	// ResourcesStatusCol is the column where the status indicator starts
	// in the resources text output.
	ResourcesStatusCol = 52
)

// Resource threshold constants for health evaluation.
const (
	// ThresholdMemoryWarnPct is the memory usage percentage that triggers a warning.
	ThresholdMemoryWarnPct = 80
	// ThresholdMemoryDangerPct is the memory usage percentage that triggers a danger alert.
	ThresholdMemoryDangerPct = 90
	// ThresholdSwapWarnPct is the swap usage percentage that triggers a warning.
	ThresholdSwapWarnPct = 50
	// ThresholdSwapDangerPct is the swap usage percentage that triggers a danger alert.
	ThresholdSwapDangerPct = 75
	// ThresholdDiskWarnPct is the disk usage percentage that triggers a warning.
	ThresholdDiskWarnPct = 85
	// ThresholdDiskDangerPct is the disk usage percentage that triggers a danger alert.
	ThresholdDiskDangerPct = 95
	// ThresholdLoadWarnRatio is the load-to-CPU ratio that triggers a warning.
	ThresholdLoadWarnRatio = 0.8
	// ThresholdLoadDangerRatio is the load-to-CPU ratio that triggers a danger alert.
	ThresholdLoadDangerRatio = 1.5
	// BytesPerGiB is the number of bytes in one gibibyte.
	BytesPerGiB = 1 << 30
)

// Ceremony configuration.
const (
	// CeremonyThrottleID is the state file name for daily throttle of ceremony checks.
	CeremonyThrottleID = "ceremony-reminded"
	// CeremonyJournalLookback is the number of recent journal files to scan for ceremony usage.
	CeremonyJournalLookback = 3
	// CeremonyRememberCmd is the command name scanned in journals for /ctx-remember usage.
	CeremonyRememberCmd = "ctx-remember"
	// CeremonyWrapUpCmd is the command name scanned in journals for /ctx-wrap-up usage.
	CeremonyWrapUpCmd = "ctx-wrap-up"
)

// Check-journal configuration.
const (
	// CheckJournalThrottleID is the state file name for daily throttle of journal checks.
	CheckJournalThrottleID = "journal-reminded"
	// CheckJournalClaudeProjectsSubdir is the relative path under $HOME to
	// the Claude Code projects directory scanned for unexported sessions.
	CheckJournalClaudeProjectsSubdir = ".claude/projects"
)

// Check-task-completion configuration.
const (
	// TaskNudgePrefix is the state file prefix for per-session
	// task completion nudge counters.
	TaskNudgePrefix = "task-nudge-"
)

// Check-resources configuration.
const (
	// CheckResourcesDangerMarker is the unicode cross marker for danger alerts.
	CheckResourcesDangerMarker = "\u2716 "
)

// Check-persistence configuration.
const (
	// PersistenceNudgePrefix is the state file prefix for per-session
	// persistence nudge counters.
	PersistenceNudgePrefix = "persistence-nudge-"
	// PersistenceEarlyMin is the minimum prompt count before nudging begins.
	PersistenceEarlyMin = 11
	// PersistenceEarlyMax is the upper bound for the early nudge window.
	PersistenceEarlyMax = 25
	// PersistenceEarlyInterval is the number of prompts between nudges
	// during the early window (prompts 11-25).
	PersistenceEarlyInterval = 20
	// PersistenceLateInterval is the number of prompts between nudges
	// after the early window (prompts 25+).
	PersistenceLateInterval = 15
	// PersistenceLogFile is the log filename for persistence check events.
	PersistenceLogFile = "check-persistence.log"
	// PersistenceKeyCount is the state file key for prompt count.
	PersistenceKeyCount = "count"
	// PersistenceKeyLastNudge is the state file key for last nudge prompt number.
	PersistenceKeyLastNudge = "last_nudge"
	// PersistenceKeyLastMtime is the state file key for last modification time.
	PersistenceKeyLastMtime = "last_mtime"
)

// Check-version configuration.
const (
	// VersionThrottleID is the state file name for daily throttle of version checks.
	VersionThrottleID = "version-checked"
	// VersionDevBuild is the version string used for development builds.
	VersionDevBuild = "dev"
)

// Context-size event names.
const (
	// EventSuppressed is the event name for suppressed prompts.
	EventSuppressed = "suppressed"
	// EventSilent is the event name for silent (no-action) prompts.
	EventSilent = "silent"
	// EventCheckpoint is the event name for context checkpoint emissions.
	EventCheckpoint = "checkpoint"
	// EventWindowWarning is the event name for context window warning emissions.
	EventWindowWarning = "window-warning"
)

// Knowledge hook configuration.
const (
	// KnowledgeThrottleID is the state file name for daily throttle of knowledge checks.
	KnowledgeThrottleID = "check-knowledge"
)

// Map staleness hook configuration.
const (
	// MapStaleDays is the threshold in days before a map refresh is considered stale.
	MapStaleDays = 30
	// MapStalenessThrottleID is the state file name for daily throttle of map staleness checks.
	MapStalenessThrottleID = "check-map-staleness"
)

// Wrap-up marker configuration.
const (
	// WrappedUpMarker is the state file name for the wrap-up suppression marker.
	WrappedUpMarker = "ctx-wrapped-up"
	// WrappedUpContent is the content written to the wrap-up marker file.
	WrappedUpContent = "wrapped-up"
)

// Tool identifiers for session parsers.
const (
	// ToolClaudeCode is the tool identifier for Claude Code sessions.
	ToolClaudeCode = "claude-code"
	// ToolMarkdown is the tool identifier for Markdown session files.
	ToolMarkdown = "markdown"
)

// Claude Code integration file names.
const (
	// FileClaudeMd is the Claude Code configuration file in the project root.
	FileClaudeMd = "CLAUDE.md"
	// FilePromptMd is the session prompt file in the project root.
	FilePromptMd = "PROMPT.md"
	// FileImplementationPlan is the implementation plan file in the project root.
	FileImplementationPlan = "IMPLEMENTATION_PLAN.md"
	// FileSettings is the Claude Code local settings file.
	FileSettings = ".claude/settings.local.json"
	// FileSettingsGolden is the golden image of the Claude Code settings.
	FileSettingsGolden = ".claude/settings.golden.json"
	// FileMakefileCtx is the ctx-owned Makefile include for project root.
	FileMakefileCtx = "Makefile.ctx"

	// FileGlobalSettings is the Claude Code global settings file.
	// Located at ~/.claude/settings.json (not the project-local one).
	FileGlobalSettings = "settings.json"
	// FileInstalledPlugins is the Claude Code installed plugins registry.
	// Located at ~/.claude/plugins/installed_plugins.json.
	FileInstalledPlugins = "plugins/installed_plugins.json"

	// PluginID is the ctx plugin identifier in Claude Code.
	PluginID = "ctx@activememory-ctx"
)

// Context file name constants for .context/ directory.
const (
	// FileConstitution contains inviolable rules for agents.
	FileConstitution = "CONSTITUTION.md"
	// FileTask contains current work items and their status.
	FileTask = "TASKS.md"
	// FileConvention contains code patterns and standards.
	FileConvention = "CONVENTIONS.md"
	// FileArchitecture contains system structure documentation.
	FileArchitecture = "ARCHITECTURE.md"
	// FileDecision contains architectural decisions with rationale.
	FileDecision = "DECISIONS.md"
	// FileLearning contains gotchas, tips, and lessons learned.
	FileLearning = "LEARNINGS.md"
	// FileGlossary contains domain terms and definitions.
	FileGlossary = "GLOSSARY.md"
	// FileAgentPlaybook contains the meta-instructions for using the
	// context system.
	FileAgentPlaybook = "AGENT_PLAYBOOK.md"
	// FileDependency contains project dependency documentation.
	FileDependency = "DEPENDENCIES.md"
)

// Journal state file.
const (
	// FileJournalState is the processing state file in .context/journal/.
	FileJournalState = ".state.json"
)

// Journal processing stage names.
const (
	// StageExported marks a journal entry as exported from Claude Code.
	StageExported = "exported"
	// StageEnriched marks a journal entry as enriched with metadata.
	StageEnriched = "enriched"
	// StageNormalized marks a journal entry as normalized for rendering.
	StageNormalized = "normalized"
	// StageFencesVerified marks a journal entry as having verified code fences.
	StageFencesVerified = "fences_verified"
	// StageLocked marks a journal entry as locked (read-only).
	StageLocked = "locked"
)

// Architecture mapping file constants for .context/ directory.
const (
	// FileDetailedDesign is the deep per-module architecture reference.
	FileDetailedDesign = "DETAILED_DESIGN.md"
	// FileMapTracking is the architecture mapping coverage state file.
	FileMapTracking = "map-tracking.json"
)

// Scratchpad file constants for .context/ directory.
const (
	// FileScratchpadEnc is the encrypted scratchpad file.
	FileScratchpadEnc = "scratchpad.enc"
	// FileScratchpadMd is the plaintext scratchpad file.
	FileScratchpadMd = "scratchpad.md"
	// FileContextKey is the context encryption key file.
	FileContextKey = ".ctx.key"
	// FileNotifyEnc is the encrypted webhook URL file.
	FileNotifyEnc = ".notify.enc"
)

// Scratchpad blob constants.
const (
	// BlobSep separates the label from the base64-encoded file content.
	BlobSep = ":::"
	// MaxBlobSize is the maximum file size (pre-encoding) allowed for blob entries.
	MaxBlobSize = 64 * 1024
	// BlobTag is the display tag appended to blob labels.
	BlobTag = " [BLOB]"
)

// Reminder file constants for .context/ directory.
const (
	// FileReminders is the session-scoped reminders file.
	FileReminders = "reminders.json"
)

// Doctor check name constants — used as Result.Name values.
const (
	// DoctorCheckContextInit identifies the context initialization check.
	DoctorCheckContextInit = "context_initialized"
	// DoctorCheckRequiredFiles identifies the required files check.
	DoctorCheckRequiredFiles = "required_files"
	// DoctorCheckCtxrcValidation identifies the .ctxrc validation check.
	DoctorCheckCtxrcValidation = "ctxrc_validation"
	// DoctorCheckDrift identifies the drift detection check.
	DoctorCheckDrift = "drift"
	// DoctorCheckPluginInstalled identifies the plugin installation check.
	DoctorCheckPluginInstalled = "plugin_installed"
	// DoctorCheckPluginEnabledGlobal identifies the global plugin enablement check.
	DoctorCheckPluginEnabledGlobal = "plugin_enabled_global"
	// DoctorCheckPluginEnabledLocal identifies the local plugin enablement check.
	DoctorCheckPluginEnabledLocal = "plugin_enabled_local"
	// DoctorCheckPluginEnabled identifies the plugin enablement check (when neither scope is active).
	DoctorCheckPluginEnabled = "plugin_enabled"
	// DoctorCheckEventLogging identifies the event logging check.
	DoctorCheckEventLogging = "event_logging"
	// DoctorCheckWebhook identifies the webhook configuration check.
	DoctorCheckWebhook = "webhook"
	// DoctorCheckReminders identifies the pending reminders check.
	DoctorCheckReminders = "reminders"
	// DoctorCheckTaskCompletion identifies the task completion check.
	DoctorCheckTaskCompletion = "task_completion"
	// DoctorCheckContextSize identifies the context token size check.
	DoctorCheckContextSize = "context_size"
	// DoctorCheckContextFilePrefix is the prefix for per-file context size results.
	DoctorCheckContextFilePrefix = "context_file_"
	// DoctorCheckRecentEvents identifies the recent event log check.
	DoctorCheckRecentEvents = "recent_events"
	// DoctorCheckResourceMemory identifies the memory resource check.
	DoctorCheckResourceMemory = "resource_memory"
	// DoctorCheckResourceSwap identifies the swap resource check.
	DoctorCheckResourceSwap = "resource_swap"
	// DoctorCheckResourceDisk identifies the disk resource check.
	DoctorCheckResourceDisk = "resource_disk"
	// DoctorCheckResourceLoad identifies the load resource check.
	DoctorCheckResourceLoad = "resource_load"
)

// Doctor category constants — used as Result.Category values.
const (
	// DoctorCategoryStructure groups context directory and file checks.
	DoctorCategoryStructure = "Structure"
	// DoctorCategoryQuality groups drift and content quality checks.
	DoctorCategoryQuality = "Quality"
	// DoctorCategoryPlugin groups plugin installation and enablement checks.
	DoctorCategoryPlugin = "Plugin"
	// DoctorCategoryHooks groups hook configuration checks.
	DoctorCategoryHooks = "Hooks"
	// DoctorCategoryState groups runtime state checks.
	DoctorCategoryState = "State"
	// DoctorCategorySize groups token size and budget checks.
	DoctorCategorySize = "Size"
	// DoctorCategoryResources groups system resource checks.
	DoctorCategoryResources = "Resources"
	// DoctorCategoryEvents groups event log checks.
	DoctorCategoryEvents = "Events"
)

// Memory bridge file constants for .context/memory/ directory.
const (
	// FileMemorySource is the Claude Code auto memory filename.
	FileMemorySource = "MEMORY.md"
	// FileMemoryMirror is the raw copy of Claude Code's MEMORY.md.
	FileMemoryMirror = "mirror.md"
	// FileMemoryState is the sync/import tracking state file.
	FileMemoryState = "memory-import.json"
)

// PathMemoryMirror is the relative path from the project root to the
// memory mirror file. Constructed from directory and file constants.
var PathMemoryMirror = filepath.Join(dir.Context, dir.Memory, FileMemoryMirror)

// Event log constants for .context/state/ directory.
const (
	// FileEventLog is the current event log file.
	FileEventLog = "events.jsonl"
	// FileEventLogPrev is the rotated (previous) event log file.
	FileEventLogPrev = "events.1.jsonl"
	// EventLogMaxBytes is the size threshold for log rotation (1MB).
	EventLogMaxBytes = 1 << 20
	// LogMaxBytes is the size threshold for hook log rotation (1MB).
	LogMaxBytes = 1 << 20
)

// FileType maps short names to actual file names.
var FileType = map[string]string{
	entry.Decision:   FileDecision,
	entry.Task:       FileTask,
	entry.Learning:   FileLearning,
	entry.Convention: FileConvention,
}

// FilesRequired lists the essential context files that must be present.
//
// These are the files created with `ctx init --minimal` and checked by
// drift detection for missing files.
var FilesRequired = []string{
	FileConstitution,
	FileTask,
	FileDecision,
}

// FileReadOrder defines the priority order for reading context files.
//
// The order follows a logical progression for AI agents:
//
//  1. CONSTITUTION — Inviolable rules. Must be loaded first so the agent
//     knows what it cannot do before attempting anything.
//
//  2. TASKS — Current work items. What the agent should focus on.
//
//  3. CONVENTIONS — How to write code. Patterns and standards to follow.
//
//  4. ARCHITECTURE — System structure. Understanding of components and
//     boundaries before making changes.
//
//  5. DECISIONS — Historical context. Why things are the way they are,
//     to avoid re-debating settled decisions.
//
//  6. LEARNINGS — Gotchas and tips. Lessons from past work that inform
//     current implementation.
//
//  7. GLOSSARY — Reference material. Domain terms and abbreviations for
//     lookup as needed.
//
//  8. AGENT_PLAYBOOK — Meta instructions. How to use this context system.
//     Loaded last because it's about the system itself, not the work.
//     The agent should understand the content before the operating manual.
var FileReadOrder = []string{
	FileConstitution,
	FileTask,
	FileConvention,
	FileArchitecture,
	FileDecision,
	FileLearning,
	FileGlossary,
	FileAgentPlaybook,
}

// Packages maps dependency manifest files to their descriptions.
//
// Nudge box drawing constants.
const (
	// BoxTop is the top-left corner of a nudge box.
	BoxTop = "┌─ "
	// BoxLinePrefix is the left border prefix for nudge box content lines.
	BoxLinePrefix = "│ "
	// BoxBottom is the bottom border of a nudge box.
	BoxBottom = "└──────────────────────────────────────────────────"
	// NudgeBoxWidth is the inner character width of the nudge box border.
	NudgeBoxWidth = 51
)

// Session and template constants.
const (
	// SessionUnknown is the fallback session ID when input lacks one.
	SessionUnknown = "unknown"
	// TemplateName is the name used for Go text/template instances.
	TemplateName = "msg"
)

// Used by sync to detect projects and suggest dependency documentation.
var Packages = map[string]string{
	"package.json":     "Node.js dependencies",
	"go.mod":           "Go module dependencies",
	"Cargo.toml":       "Rust dependencies",
	"requirements.txt": "Python dependencies",
	"Gemfile":          "Ruby dependencies",
}
