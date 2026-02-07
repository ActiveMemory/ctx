//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"path"

	"github.com/ActiveMemory/ctx/internal/config"
)

func NewHook(hookType HookType, cmd string) Hook {
	return Hook{
		Type:    hookType,
		Command: cmd,
	}
}

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
		hooksDir = path.Join(projectDir, config.DirClaudeHooks)
	}

	return HookConfig{
		PreToolUse:       preToolUserHookMatcher(hooksDir),
		UserPromptSubmit: userPromptSubmitHookMatcher(hooksDir),
		SessionEnd:       sessionEndHookMatcher(hooksDir),
	}
}
