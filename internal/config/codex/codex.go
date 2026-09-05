//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

// Codex CLI binary and home directory.
const (
	// Binary is the Codex CLI binary name, resolved via
	// exec.LookPath to detect whether Codex is installed.
	Binary = "codex"

	// DirHome is the default Codex home directory name under
	// $HOME (overridden by [EnvHome]).
	DirHome = ".codex"
	// EnvHome is the environment variable that relocates the
	// Codex home directory.
	EnvHome = "CODEX_HOME"

	// DirSessions is the rollout transcript directory under the
	// Codex home (sessions/YYYY/MM/DD/rollout-*.jsonl).
	DirSessions = "sessions"
	// RolloutPrefix is the filename prefix of Codex rollout
	// transcripts.
	RolloutPrefix = "rollout-"

	// DirPlugins is the plugins directory under the Codex home.
	DirPlugins = "plugins"
	// DirPluginCache is the installed-plugin cache under
	// [DirPlugins] (cache/<marketplace>/<plugin>/<version>/).
	DirPluginCache = "cache"
	// PluginVersionLocal is the <version> segment Codex uses for
	// plugins installed from a local marketplace.
	PluginVersionLocal = "local"
)

// Project-local Codex layout.
const (
	// Dir is the project-local Codex config directory.
	Dir = ".codex"
	// FileHooksJSON is the hooks manifest file name (project
	// `.codex/hooks.json`, plugin `hooks/hooks.json`).
	FileHooksJSON = "hooks.json"
	// FileConfigTOML is the Codex config file name (user
	// `~/.codex/config.toml`, project `.codex/config.toml`).
	FileConfigTOML = "config.toml"
	// DirHooks is the hooks subdirectory inside a plugin root.
	DirHooks = "hooks"

	// DirAgents is the cross-agent directory Codex scans for
	// repo-scoped skills and marketplaces.
	DirAgents = ".agents"
	// DirMarketplacePlugins is the marketplace subdirectory under
	// [DirAgents] (`.agents/plugins/marketplace.json`).
	DirMarketplacePlugins = "plugins"
	// FileMarketplaceJSON is the marketplace catalog file name.
	FileMarketplaceJSON = "marketplace.json"

	// FileMCPJSON is the plugin-bundled MCP server map file name.
	FileMCPJSON = ".mcp.json"
	// DirPluginManifest is the manifest directory inside a plugin
	// root (`.codex-plugin/`). Its presence in an installed cache
	// copy identifies the Codex plugin variant; a cache holding
	// the legacy Claude Code variant has `.claude-plugin/` instead.
	DirPluginManifest = ".codex-plugin"
)

// Plugin identity.
const (
	// PluginName is the ctx plugin name in the Codex manifest.
	PluginName = "ctx"
	// MarketplaceID is the repo marketplace name
	// (`.agents/plugins/marketplace.json` → `name`).
	MarketplaceID = "activememory-ctx"
	// PluginID is the `<plugin>@<marketplace>` identifier Codex
	// uses in `codex plugin add` and in config.toml.
	PluginID = PluginName + "@" + MarketplaceID
)

// config.toml tokens.
//
// ctx never round-trips Codex's TOML through a parser: the user
// owns that file (comments, ordering). The deployer appends a
// table when its header is absent and scans for headers when
// detecting state; these constants are the exact header lines
// it writes and looks for.
const (
	// TOMLHeaderMCPCtx is the table header for the ctx MCP server.
	TOMLHeaderMCPCtx = "[mcp_servers.ctx]"
	// TOMLHeaderPluginCtx is the table header Codex writes when the
	// ctx plugin is installed/enabled.
	TOMLHeaderPluginCtx = `[plugins."` + PluginID + `"]`
	// TOMLKeyEnabled is the enabled flag key inside a plugin table.
	TOMLKeyEnabled = "enabled"
	// TOMLKeyCommand is the MCP server command key.
	TOMLKeyCommand = "command"
	// TOMLKeyArgs is the MCP server args key.
	TOMLKeyArgs = "args"
	// TOMLTrue is the TOML boolean literal for true.
	TOMLTrue = "true"
	// TOMLBracketOpen opens an inline array and, at the start of
	// a trimmed line, a table header.
	TOMLBracketOpen = "["
	// TOMLBracketClose closes an inline array.
	TOMLBracketClose = "]"
	// TOMLDot separates key segments in a table header.
	TOMLDot = "."
	// TOMLComment starts a TOML comment.
	TOMLComment = "#"
	// TOMLQuoteBasic is the basic-string quote for a key segment.
	TOMLQuoteBasic = `"`
	// TOMLQuoteLiteral is the literal-string quote for a key segment.
	TOMLQuoteLiteral = "'"
)

// Skill frontmatter tokens.
const (
	// FrontmatterKeyName is the SKILL.md frontmatter key whose
	// value must equal the skill directory name for a file to
	// count as ctx-managed.
	FrontmatterKeyName = "name"
)

// Hook manifest tokens.
const (
	// KeyHooks is the top-level key of hooks.json.
	KeyHooks = "hooks"
	// KeyDescription is the optional top-level metadata key.
	KeyDescription = "description"
	// KeyHandlers is the handler-array key inside a matcher group
	// (also spelled "hooks" in the manifest).
	KeyHandlers = "hooks"
	// KeyCommand is the handler command key.
	KeyCommand = "command"
	// KeyType is the handler type key.
	KeyType = "type"
	// HandlerTypeCommand is the only handler type Codex runs today.
	HandlerTypeCommand = "command"

	// HookAnchor is the git-root anchor every ctx hook command
	// starts with. Codex runs hooks with the session cwd (which
	// may be a subdirectory), and ctx is CWD-anchored, so the
	// command must `cd` to the project root first.
	HookAnchor = `cd "$(git rev-parse --show-toplevel 2>/dev/null || pwd)" && `

	// LegacyHookAnchor is the anchor shipped by earlier builds
	// (no fallback when the cwd is outside a git repo — Codex can
	// run hooks from non-repo directories, where the bare form
	// exits 1 before ctx starts). Recognized so merges migrate
	// previously deployed groups instead of duplicating them.
	LegacyHookAnchor = `cd "$(git rev-parse --show-toplevel)" && `
	// LegacyHookCommandPrefixGuardless is the tolerant anchor
	// without the ctx-absent guard — the shape deployed between
	// the anchor fix and the guard fix.
	LegacyHookCommandPrefixGuardless = HookAnchor + "ctx "
	// LegacyHookCommandPrefix is [LegacyHookAnchor] followed by a
	// ctx invocation — the ctx-managed command shape the earliest
	// builds deployed.
	LegacyHookCommandPrefix = LegacyHookAnchor + "ctx "

	// HookGuard exits a hook silently when the ctx binary is not
	// on PATH: ctx is an optional companion, and a collaborator
	// without it must not see a wall of exit-127 hook failures.
	HookGuard = `command -v ctx >/dev/null 2>&1 || exit 0; `

	// HookPrologue is the full prefix of every shipped hook
	// command: the ctx-absent guard followed by the git-root
	// anchor.
	HookPrologue = HookGuard + HookAnchor

	// HookCommandPrefix is [HookPrologue] followed by a ctx
	// binary invocation. Ownership checks require this whole
	// prefix so user hooks that merely copy the anchor idiom are
	// never classified as ctx-managed.
	HookCommandPrefix = HookPrologue + "ctx "

	// SessionEndTimeoutMax is the maximum timeout (seconds) Codex
	// allows for SessionEnd hooks.
	SessionEndTimeoutMax = 3
)

// Lifecycle hook event names Codex supports.
const (
	EventSessionStart      = "SessionStart"
	EventSessionEnd        = "SessionEnd"
	EventUserPromptSubmit  = "UserPromptSubmit"
	EventPreToolUse        = "PreToolUse"
	EventPostToolUse       = "PostToolUse"
	EventPermissionRequest = "PermissionRequest"
	EventPreCompact        = "PreCompact"
	EventPostCompact       = "PostCompact"
	EventSubagentStart     = "SubagentStart"
	EventSubagentStop      = "SubagentStop"
	EventStop              = "Stop"
)

// Events lists every lifecycle event Codex supports; the manifest
// guard rejects any other key.
var Events = []string{
	EventSessionStart,
	EventSessionEnd,
	EventUserPromptSubmit,
	EventPreToolUse,
	EventPostToolUse,
	EventPermissionRequest,
	EventPreCompact,
	EventPostCompact,
	EventSubagentStart,
	EventSubagentStop,
	EventStop,
}

// Tool names as Codex reports them to hooks.
const (
	// ToolBash is the canonical hook name for shell commands.
	ToolBash = "Bash"
	// ToolApplyPatch is the canonical hook name for file edits.
	ToolApplyPatch = "apply_patch"
	// ToolUpdatePlan is Codex's planning tool (the analogue of
	// Claude Code's EnterPlanMode).
	ToolUpdatePlan = "update_plan"
)

// Rollout transcript tokens (sessions/YYYY/MM/DD/rollout-*.jsonl).
const (
	// LineTypeSessionMeta is the first line of every rollout.
	LineTypeSessionMeta = "session_meta"
	// LineTypeResponseItem carries messages, tool calls, and
	// tool outputs.
	LineTypeResponseItem = "response_item"
	// LineTypeEventMsg carries UI events (token counts, task
	// lifecycle, duplicated message text).
	LineTypeEventMsg = "event_msg"
	// LineTypeTurnContext carries per-turn model/cwd context.
	LineTypeTurnContext = "turn_context"

	// ItemTypeMessage is a response_item message.
	ItemTypeMessage = "message"
	// ItemTypeFunctionCall is a response_item function-tool
	// invocation (name + JSON arguments).
	ItemTypeFunctionCall = "function_call"
	// ItemTypeFunctionCallOutput is a response_item function-tool
	// result.
	ItemTypeFunctionCallOutput = "function_call_output"
	// ItemTypeCustomToolCall is a response_item custom-tool
	// invocation (name + free-form input, e.g. code-mode `exec`).
	ItemTypeCustomToolCall = "custom_tool_call"
	// ItemTypeCustomToolCallOutput is a response_item custom-tool
	// result.
	ItemTypeCustomToolCallOutput = "custom_tool_call_output"
	// ItemTypeLocalShellCall is a response_item shell invocation
	// used by older Codex releases.
	ItemTypeLocalShellCall = "local_shell_call"

	// ContentInputText is a user content part.
	ContentInputText = "input_text"
	// ContentOutputText is an assistant content part.
	ContentOutputText = "output_text"

	// RoleUser is the user message role.
	RoleUser = "user"
	// RoleAssistant is the assistant message role.
	RoleAssistant = "assistant"

	// EventTokenCount is the event_msg type carrying token usage.
	EventTokenCount = "token_count"
	// EventItemCompleted is the event_msg type that mirrors a
	// finished item; its CommandExecution items carry exit codes.
	EventItemCompleted = "item_completed"
	// ItemCommandExecution is the item_completed item type for a
	// shell command.
	ItemCommandExecution = "CommandExecution"
)

// InjectedUserPrefixes are the opening markers of user-role items
// Codex injects itself (environment, instructions, permissions,
// AGENTS.md payloads). They are not user prose and are dropped from
// imported sessions.
var InjectedUserPrefixes = []string{
	"# AGENTS.md instructions for ",
	"<environment_context>",
	"<user_instructions>",
	"<permissions instructions>",
	"<recommended_plugins>",
	"<skills_instructions>",
	"<apps_instructions>",
	"<plugins_instructions>",
	"<multi_agent_mode>",
	"<turn_aborted>",
	"<skill>",
	"<app_instructions>",
}
