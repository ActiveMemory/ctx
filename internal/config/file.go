//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

// File permission constants.
const (
	// PermFile is the standard permission for regular files (owner rw, others r).
	PermFile = 0644
	// PermExec is the standard permission for directories and executable files.
	PermExec = 0755
)

// File extension constants.
const (
	// ExtMarkdown is the Markdown file extension.
	ExtMarkdown = ".md"
	// ExtJSONL is the JSON Lines file extension.
	ExtJSONL = ".jsonl"
)

// Common filenames.
const (
	// FilenameReadme is the standard README filename.
	FilenameReadme = "README.md"
	// FilenameIndex is the standard index filename for generated sites.
	FilenameIndex = "index.md"
)

// Journal site configuration.
const (
	// FileZensicalToml is the zensical site configuration filename.
	FileZensicalToml = "zensical.toml"
	// BinZensical is the zensical binary name.
	BinZensical = "zensical"
)

// Generated file suffixes.
const (
	// SuffixPreCompact is the suffix for pre-compact session snapshots.
	SuffixPreCompact = "-pre-compact.md"
	// SuffixSummary is the suffix for summary files that accompany transcripts.
	SuffixSummary = "-summary.md"
)

// Session defaults.
const (
	// DefaultSessionTopic is the fallback topic when none is provided.
	DefaultSessionTopic = "manual-save"
	// DefaultSessionFilename is the fallback filename component when
	// sanitization produces an empty string.
	DefaultSessionFilename = "session"
)

// Runtime configuration constants.
const (
	// FileContextRC is the optional runtime configuration file.
	FileContextRC = ".contextrc"
)

// Environment configuration.
const (
	// EnvCtxDir is the environment variable for overriding the context directory.
	EnvCtxDir = "CTX_DIR"
	// EnvCtxTokenBudget is the environment variable for overriding the token budget.
	EnvCtxTokenBudget = "CTX_TOKEN_BUDGET" //nolint:gosec // G101: env var name, not a credential
	// EnvCtxSessionStart is the environment variable for session start correlation.
	EnvCtxSessionStart = "CTX_SESSION_START"
)

// Parser configuration.
const (
	// ParserPeekLines is the number of lines to scan when detecting file format.
	ParserPeekLines = 50
)

// Claude API content block types.
const (
	// ClaudeBlockText is a text content block.
	ClaudeBlockText = "text"
	// ClaudeBlockThinking is an extended thinking content block.
	ClaudeBlockThinking = "thinking"
	// ClaudeBlockToolUse is a tool invocation block.
	ClaudeBlockToolUse = "tool_use"
	// ClaudeBlockToolResult is a tool execution result block.
	ClaudeBlockToolResult = "tool_result"
)

// Claude API content block field keys.
const (
	// ClaudeFieldType is the block type discriminator key.
	ClaudeFieldType = "type"
	// ClaudeFieldText is the text content key.
	ClaudeFieldText = "text"
	// ClaudeFieldThinking is the thinking content key.
	ClaudeFieldThinking = "thinking"
	// ClaudeFieldName is the tool name key.
	ClaudeFieldName = "name"
	// ClaudeFieldInput is the tool input parameters key.
	ClaudeFieldInput = "input"
	// ClaudeFieldContent is the tool result content key.
	ClaudeFieldContent = "content"
)

// Claude API message roles.
const (
	// RoleUser is a user message.
	RoleUser = "user"
	// RoleAssistant is an assistant message.
	RoleAssistant = "assistant"
)

// Claude Code integration file names.
const (
	// FileAutoSave is the hook script for auto-saving sessions.
	FileAutoSave = "auto-save-session.sh"
	// FileBlockNonPathScript is the hook script that blocks non-PATH ctx
	// invocations.
	FileBlockNonPathScript = "block-non-path-ctx.sh"
	// FilePromptCoach is the hook script for prompt quality feedback.
	FilePromptCoach = "prompt-coach.sh"
	// FileCheckContextSize is the hook script for context size checkpoints.
	FileCheckContextSize = "check-context-size.sh"
	// FileClaudeMd is the Claude Code configuration file in the project root.
	FileClaudeMd = "CLAUDE.md"
	// FilePromptMd is the session prompt file in the project root.
	FilePromptMd = "PROMPT.md"
	// FileImplementationPlan is the implementation plan file in the project root.
	FileImplementationPlan = "IMPLEMENTATION_PLAN.md"
	// FileSettings is the Claude Code local settings file.
	FileSettings = ".claude/settings.local.json"
	// FileContextWatch is the context monitoring tool script.
	FileContextWatch = "context-watch.sh"
	// FileMakefileCtx is the ctx-owned Makefile include for project root.
	FileMakefileCtx = "Makefile.ctx"
	// CmdAutoloadContext is the inline command for the PreToolUse hook
	// that autoloads the context packet on every tool use. The cooldown
	// inside ctx agent prevents repetitive output.
	CmdAutoloadContext = "ctx agent --budget 4000 --session $PPID 2>/dev/null || true"
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

// FileType maps short names to actual file names.
var FileType = map[string]string{
	EntryDecision:   FileDecision,
	EntryTask:       FileTask,
	EntryLearning:   FileLearning,
	EntryConvention: FileConvention,
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
// Used by sync to detect projects and suggest dependency documentation.
var Packages = map[string]string{
	"package.json":     "Node.js dependencies",
	"go.mod":           "Go module dependencies",
	"Cargo.toml":       "Rust dependencies",
	"requirements.txt": "Python dependencies",
	"Gemfile":          "Ruby dependencies",
}

// DefaultClaudePermissions lists the default permissions for ctx commands.
//
// These permissions allow Claude Code to run ctx CLI commands without
// prompting for approval. All ctx subcommands are pre-approved.
var DefaultClaudePermissions = []string{
	"Bash(ctx status:*)",
	"Bash(ctx agent:*)",
	"Bash(ctx add:*)",
	"Bash(ctx session:*)",
	"Bash(ctx tasks:*)",
	"Bash(ctx loop:*)",
}
