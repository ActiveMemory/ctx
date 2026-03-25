//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

// Hook name constants: used for LoadMessage, NewTemplateRef, notify.Send,
// and log.Append to avoid magic strings.
const (
	// BlockDangerousCommands is the hook name for blocking dangerous commands.
	BlockDangerousCommands = "block-dangerous-commands"
	// BlockNonPathCtx is the hook name for blocking non-PATH ctx invocations.
	BlockNonPathCtx = "block-non-path-ctx"
	// CheckBackupAge is the hook name for backup staleness checks.
	CheckBackupAge = "check-backup-age"
	// CheckCeremonies is the hook name for ceremony usage checks.
	CheckCeremonies = "check-ceremonies"
	// CheckContextSize is the hook name for context window size checks.
	CheckContextSize = "check-context-size"
	// CheckFreshness is the hook name for technology constant freshness checks.
	CheckFreshness = "check-freshness"
	// CheckJournal is the hook name for journal health checks.
	CheckJournal = "check-journal"
	// CheckKnowledge is the hook name for knowledge file health checks.
	CheckKnowledge = "check-knowledge"
	// CheckMapStaleness is the hook name for architecture map staleness checks.
	CheckMapStaleness = "check-map-staleness"
	// CheckMemoryDrift is the hook name for memory drift checks.
	CheckMemoryDrift = "check-memory-drift"
	// CheckPersistence is the hook name for context persistence nudges.
	CheckPersistence = "check-persistence"
	// CheckReminders is the hook name for session reminder checks.
	CheckReminders = "check-reminders"
	// CheckResources is the hook name for resource usage checks.
	CheckResources = "check-resources"
	// CheckTaskCompletion is the hook name for task completion nudges.
	CheckTaskCompletion = "check-task-completion"
	// CheckVersion is the hook name for version mismatch checks.
	CheckVersion = "check-version"
	// Heartbeat is the hook name for session heartbeat events.
	Heartbeat = "heartbeat"
	// PostCommit is the hook name for post-commit nudges.
	PostCommit = "post-commit"
	// QAReminder is the hook name for QA reminder gates.
	QAReminder = "qa-reminder"
	// SpecsNudge is the hook name for specs directory nudges.
	SpecsNudge = "specs-nudge"
	// VersionDrift is the hook name for version drift nudges.
	VersionDrift = "version-drift"
)

// Supported integration tool names for ctx hook command.
const (
	ToolAgents     = "agents"
	ToolAider      = "aider"
	ToolClaude     = "claude"
	ToolClaudeCode = "claude-code"
	ToolCopilot    = "copilot"
	ToolCopilotCLI = "copilot-cli"
	ToolCursor     = "cursor"
	ToolWindsurf   = "windsurf"
)

// Copilot integration paths.
const (
	DirGitHub               = ".github"
	DirGitHubAgents         = "agents"
	DirGitHubHooks          = "hooks"
	DirGitHubHooksScripts   = "scripts"
	DirGitHubInstructions   = "instructions"
	DirGitHubSkills         = "skills"
	FileAgentsMd            = "AGENTS.md"
	FileAgentsCtxMd         = "ctx.md"
	FileCopilotInstructions = "copilot-instructions.md"
	FileCopilotCLIHooksJSON = "ctx-hooks.json"
	FileInstructionsCtxMd   = "context.instructions.md"
	FileSKILLMd             = "SKILL.md"
)

// Copilot CLI home directory and MCP config.
const (
	// DirCopilotHome is the default Copilot CLI config directory name.
	DirCopilotHome = ".copilot"
	// EnvCopilotHome is the environment variable to override the config dir.
	EnvCopilotHome = "COPILOT_HOME"
	// FileMCPConfigJSON is the MCP server configuration file name.
	FileMCPConfigJSON = "mcp-config.json"
)

// Prefixes
const (
	// StdinReadTimeout is the maximum time to wait for hook JSON on stdin
	// before returning a zero-value input.
	StdinReadTimeout = 2

	// PrefixMemoryDriftThrottle is the state file prefix for per-session
	// memory drift nudge tombstones.
	PrefixMemoryDriftThrottle = "memory-drift-nudged-"
	// PrefixPauseMarker is the state file prefix for session pause markers.
	PrefixPauseMarker = "ctx-paused-"
	// LabelPaused is the short status label emitted while hooks are paused.
	LabelPaused = "ctx:paused"
)

// Hook event names (Claude Code hook lifecycle stages).
const (
	// EventPreToolUse is the hook event for pre-tool-use hooks.
	EventPreToolUse = "PreToolUse"
	// EventPostToolUse is the hook event for post-tool-use hooks.
	EventPostToolUse = "PostToolUse"
)

// Copilot CLI hook event names (GitHub Copilot CLI lifecycle stages).
const (
	CLIEventSessionStart = "sessionStart"
	CLIEventSessionEnd   = "sessionEnd"
	CLIEventPreToolUse   = "preToolUse"
	CLIEventPostToolUse  = "postToolUse"
)
