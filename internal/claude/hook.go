//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import "fmt"

// DefaultHooks returns the default ctx hooks configuration for
// Claude Code.
//
// The returned hooks configure PreToolUse to block non-PATH ctx
// invocations and autoload context on every tool use, UserPromptSubmit
// for prompt coaching, and SessionEnd to run auto-save-session.sh for
// persisting session transcripts.
//
// Parameters:
//   - projectDir: Project root directory for absolute hook paths; if empty,
//     paths are relative (e.g., ".claude/hooks/")
//
// Returns:
//   - HookConfig: Configured hooks for PreToolUse, UserPromptSubmit, and
//     SessionEnd events
func DefaultHooks(projectDir string) HookConfig {
	hooksDir := ".claude/hooks"
	if projectDir != "" {
		hooksDir = fmt.Sprintf("%s/.claude/hooks", projectDir)
	}

	return HookConfig{
		PreToolUse: []HookMatcher{
			{
				// Block non-PATH ctx invocations (./ctx, ./dist/ctx, go run ./cmd/ctx)
				Matcher: "Bash",
				Hooks: []Hook{
					{
						Type:    "command",
						Command: fmt.Sprintf("%s/block-non-path-ctx.sh", hooksDir),
					},
				},
			},
			{
				// Autoload context on every tool use
				Matcher: ".*",
				Hooks: []Hook{
					{
						Type:    "command",
						Command: "ctx agent --budget 4000 2>/dev/null || true",
					},
				},
			},
		},
		UserPromptSubmit: []HookMatcher{
			{
				// Prompt coaching: detect anti-patterns and suggest improvements
				Hooks: []Hook{
					{
						Type:    "command",
						Command: fmt.Sprintf("%s/prompt-coach.sh", hooksDir),
					},
				},
			},
		},
		SessionEnd: []HookMatcher{
			{
				Hooks: []Hook{
					{
						Type:    "command",
						Command: fmt.Sprintf("%s/auto-save-session.sh", hooksDir),
					},
				},
			},
		},
	}
}
