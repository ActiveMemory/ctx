//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/config"
)

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
	hooksDir := config.DirClaudeHooks
	if projectDir != "" {
		hooksDir = fmt.Sprintf(
			"%s/%s", projectDir, config.DirClaudeHooks,
		)
	}

	return HookConfig{
		PreToolUse: []HookMatcher{
			{
				// Block non-PATH ctx invocations (./ctx, ./dist/ctx, go run ./cmd/ctx)
				Matcher: "Bash",
				Hooks: []Hook{
					{
						Type: "command",
						Command: fmt.Sprintf(
							"%s/%s", hooksDir, config.FileBlockNonPathScript,
						),
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
						Command: fmt.Sprintf("%s/%s", hooksDir, config.FilePromptCoach),
					},
				},
			},
		},
		SessionEnd: []HookMatcher{
			{
				Hooks: []Hook{
					{
						Type:    "command",
						Command: fmt.Sprintf("%s/%s", hooksDir, config.FileAutoSave),
					},
				},
			},
		},
	}
}
