//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package claude provides Claude Code integration templates and utilities.
//
// It embeds hook scripts, skill definitions, and configuration types
// for integrating ctx with Claude Code's settings.local.json. The embedded
// assets are installed to project directories via "ctx init --claude".
//
// Embedded assets:
//   - block-non-path-ctx.sh: Prevents non-PATH ctx invocations
//   - prompt-coach.sh: Detects prompt antipatterns and suggests improvements
//   - check-context-size.sh: Adaptive context capacity reminders
//   - check-persistence.sh: Nudges agent to persist learnings and decisions
//   - skills/*/SKILL.md: Agent skill definitions for Claude Code
//
// Example usage:
//
//	script, err := claude.BlockNonPathCtxScript()
//	if err != nil {
//	    return err
//	}
//	os.WriteFile(".claude/hooks/block-non-path-ctx.sh", script, 0755)
package claude
