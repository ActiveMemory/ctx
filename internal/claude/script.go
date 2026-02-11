//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/tpl"
)

// BlockNonPathCtxScript returns the script that blocks non-PATH ctx
// invocations.
//
// The script prevents Claude from running ctx via relative paths (./ctx,
// ./dist/ctx) or "go run", ensuring only the installed PATH version is used.
// It is installed to .claude/hooks/ during ctx init --claude.
//
// Returns:
//   - []byte: Raw bytes of the block-non-path-ctx.sh script
//   - error: Non-nil if the embedded file cannot be read
func BlockNonPathCtxScript() ([]byte, error) {
	content, err := tpl.ClaudeHookByFileName(config.FileBlockNonPathScript)
	if err != nil {
		return nil, errFileRead(config.FileBlockNonPathScript, err)
	}
	return content, nil
}

// PromptCoachScript returns the prompt coaching hook script.
//
// The script detects prompt antipatterns (e.g., "idiomatic Go") and suggests
// better alternatives (e.g., "follow project conventions"). It limits
// suggestions to 3 per pattern per session to avoid annoying the user.
// It is installed to .claude/hooks/ during ctx init --claude.
//
// Returns:
//   - []byte: Raw bytes of the prompt-coach.sh script
//   - error: Non-nil if the embedded file cannot be read
func PromptCoachScript() ([]byte, error) {
	content, err := tpl.ClaudeHookByFileName(config.FilePromptCoach)
	if err != nil {
		return nil, errFileRead(config.FilePromptCoach, err)
	}
	return content, nil
}

// CheckContextSizeScript returns the context size checkpoint hook script.
//
// The script counts prompts per session and outputs adaptive reminders to
// stderr, prompting Claude to assess remaining context capacity. Frequency
// increases as the session grows (silent for first 15 prompts, then every
// 5th, then every 3rd).
//
// Returns:
//   - []byte: Raw bytes of the check-context-size.sh script
//   - error: Non-nil if the embedded file cannot be read
func CheckContextSizeScript() ([]byte, error) {
	content, err := tpl.ClaudeHookByFileName(config.FileCheckContextSize)
	if err != nil {
		return nil, errFileRead(config.FileCheckContextSize, err)
	}
	return content, nil
}

// CheckPersistenceScript returns the persistence nudge hook script.
//
// The script counts prompts since the last .context/ file modification
// and nudges the agent to persist learnings, decisions, or task updates.
// Silent for the first 10 prompts; nudges once around prompt 20, then
// every 15th prompt if no context files have been modified.
//
// Returns:
//   - []byte: Raw bytes of the check-persistence.sh script
//   - error: Non-nil if the embedded file cannot be read
func CheckPersistenceScript() ([]byte, error) {
	content, err := tpl.ClaudeHookByFileName(config.FileCheckPersistence)
	if err != nil {
		return nil, errFileRead(config.FileCheckPersistence, err)
	}
	return content, nil
}
