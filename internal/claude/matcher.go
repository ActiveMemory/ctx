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

// preToolUserHookMatcher builds the PreToolUse hook matchers.
//
// It returns two matchers: a Bash-only matcher that blocks non-PATH ctx
// invocations (./ctx, ./dist/ctx, go run ./cmd/ctx), and a catch-all
// matcher that autoloads context via "ctx agent" on every tool use.
//
// Parameters:
//   - hooksDir: directory containing hook scripts
//
// Returns:
//   - []HookMatcher: matchers for PreToolUse lifecycle event
func preToolUserHookMatcher(hooksDir string) []HookMatcher {
	return []HookMatcher{{
		// Block non-PATH ctx invocations (./ctx, ./dist/ctx, go run ./cmd/ctx)
		Matcher: MatcherBash,
		Hooks: []Hook{NewHook(
			HookTypeCommand,
			path.Join(hooksDir, config.FileBlockNonPathScript),
		)},
	}, {
		// Autoload context on every tool use (cooldown prevents repetition)
		Matcher: MatcherAll,
		Hooks: []Hook{NewHook(
			HookTypeCommand,
			config.CmdAutoloadContext,
		)},
	}}
}

// userPromptSubmitHookMatcher builds the UserPromptSubmit hook matchers.
//
// It returns a single matcher with hooks for prompt coaching and context
// size monitoring, both triggered when the user submits a prompt.
//
// Parameters:
//   - hooksDir: directory containing hook scripts
//
// Returns:
//   - []HookMatcher: matchers for UserPromptSubmit lifecycle event
func userPromptSubmitHookMatcher(hooksDir string) []HookMatcher {
	return []HookMatcher{{
		// Prompt coaching, context monitoring, and persistence nudges
		Hooks: []Hook{
			NewHook(
				HookTypeCommand, path.Join(hooksDir, config.FilePromptCoach),
			),
			NewHook(
				HookTypeCommand, path.Join(hooksDir, config.FileCheckContextSize),
			),
			NewHook(
				HookTypeCommand, path.Join(hooksDir, config.FileCheckPersistence),
			),
		},
	}}
}
