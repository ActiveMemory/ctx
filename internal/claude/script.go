//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/tpl"
)

// AutoSaveScript returns the auto-save session script content.
//
// The script automatically saves Claude Code session transcripts when a
// session ends. It is installed to .claude/hooks/ during ctx init --claude.
//
// Returns:
//   - []byte: Raw bytes of the auto-save-session.sh script
//   - error: Non-nil if the embedded file cannot be read
func AutoSaveScript() ([]byte, error) {
	content, err := tpl.ClaudeHookByFileName(config.FileAutoSave)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", config.FileAutoSave, err)
	}
	return content, nil
}

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
		return nil, fmt.Errorf(
			"failed to read %s: %w", config.FileBlockNonPathScript, err,
		)
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
		return nil, fmt.Errorf(
			"failed to read %s: %w", config.FilePromptCoach, err,
		)
	}
	return content, nil
}
