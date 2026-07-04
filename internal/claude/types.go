//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

// HookConfig represents the hooks section of Claude Code's
// settings.local.json.
//
// Hooks are shell commands that Claude Code executes at specific lifecycle
// events. See https://docs.anthropic.com/en/docs/claude-code/hooks for details.
//
// Fields:
//   - PreToolUse: Matchers that run before each tool invocation
//   - PostToolUse: Matchers that run after a successful tool invocation
//   - UserPromptSubmit: Matchers that run when the user submits a prompt
//   - SessionEnd: Matchers that run when a session ends
type HookConfig struct {
	PreToolUse       []HookMatcher `json:"PreToolUse,omitempty"`
	PostToolUse      []HookMatcher `json:"PostToolUse,omitempty"`
	UserPromptSubmit []HookMatcher `json:"UserPromptSubmit,omitempty"`
	SessionEnd       []HookMatcher `json:"SessionEnd,omitempty"`
}

// HookType is the type identifier for a hook (e.g., "command").
type HookType string

// Matcher is a regex pattern for matching tool names in hooks.
type Matcher string

// HookMatcher associates a regex pattern with hooks to execute.
//
// For PreToolUse hooks, the Matcher pattern matches against the tool name
// (e.g., "Bash", "Read"). Use ".*" to match all tools.
//
// Fields:
//   - Matcher: Regex pattern to match; empty string matches all
//   - Hooks: Commands to execute when the pattern matches
type HookMatcher struct {
	Matcher Matcher `json:"matcher,omitempty"`
	Hooks   []Hook  `json:"hooks"`
}

// Hook represents a single hook command to execute.
//
// Fields:
//   - Type: Hook type, typically "command"
//   - Command: Shell command or script path to execute
type Hook struct {
	Type    HookType `json:"type"`
	Command string   `json:"command"`
}

// PermissionsConfig represents the permissions section of Claude Code's
// settings.local.json.
//
// Fields:
//   - Allow: List of tool patterns that are pre-approved
//     (e.g., "Bash(ctx status:*)")
//   - Deny: List of tool patterns that are always blocked
//     (e.g., "Bash(sudo *)")
type PermissionsConfig struct {
	Allow []string `json:"allow,omitempty"`
	Deny  []string `json:"deny,omitempty"`
}

// StatusLineConfig represents the statusLine section of Claude Code's
// settings.local.json.
//
// Claude Code pipes a JSON payload to the configured command on stdin
// and displays the first line of its stdout as the status line. See
// https://code.claude.com/docs/en/statusline for the payload schema.
//
// Fields:
//   - Type: Entry type, "command" for executable status lines
//   - Command: Shell command Claude Code runs on each status refresh
//   - Padding: Optional horizontal padding override
type StatusLineConfig struct {
	Type    string `json:"type,omitempty"`
	Command string `json:"command,omitempty"`
	Padding *int   `json:"padding,omitempty"`
}

// Settings represents the full Claude Code settings.local.json structure.
//
// This is used when reading or writing project-level Claude Code configuration.
//
// Fields:
//   - Hooks: Hook configuration for lifecycle events
//   - Permissions: Tool permission configuration
//   - StatusLine: Status line command configuration
type Settings struct {
	Hooks       HookConfig        `json:"hooks,omitempty"`
	Permissions PermissionsConfig `json:"permissions,omitempty"`
	StatusLine  *StatusLineConfig `json:"statusLine,omitempty"`
}
