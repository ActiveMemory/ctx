//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package claude provides Claude Code integration templates and utilities.
//
// It embeds hook scripts, slash command definitions, and configuration types
// for integrating ctx with Claude Code's settings.local.json. The embedded
// assets are installed to project directories via "ctx init --claude".
//
// Embedded assets:
//   - auto-save-session.sh: Saves session transcripts on session end
//   - block-non-path-ctx.sh: Prevents non-PATH ctx invocations
//   - tpl/commands/*.md: Slash command definitions for Claude Code
//
// Example usage:
//
//	script, err := claude.AutoSaveScript()
//	if err != nil {
//	    return err
//	}
//	os.WriteFile(".claude/hooks/auto-save-session.sh", script, 0755)
package claude
