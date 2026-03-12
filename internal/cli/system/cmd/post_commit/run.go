//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package post_commit

import (
	"os"
	"regexp"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/notify"
)

var (
	reGitCommit = regexp.MustCompile(`git\s+commit`)
	reAmend     = regexp.MustCompile(`--amend`)
)

// Run executes the post-commit hook logic.
//
// After a successful git commit (non-amend), nudges the agent to offer
// context capture (decision or learning) and to run lints/tests before
// pushing. Also checks for version drift.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	if !core.IsInitialized() {
		return nil
	}
	input, sessionID, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}

	command := input.ToolInput.Command

	// Only trigger on git commit commands
	if !reGitCommit.MatchString(command) {
		return nil
	}

	// Skip amend commits
	if reAmend.MatchString(command) {
		return nil
	}

	hook, variant := file.HookPostCommit, file.VariantNudge

	fallback := assets.TextDesc(assets.TextDescKeyPostCommitFallback)
	msg := core.LoadMessage(hook, variant, nil, fallback)
	if msg == "" {
		return nil
	}
	msg = core.AppendContextDir(msg)
	core.PrintHookContext(cmd, file.HookEventPostToolUse, msg)

	ref := notify.NewTemplateRef(hook, variant, nil)
	core.Relay(hook+": "+assets.TextDesc(assets.TextDescKeyPostCommitRelayMessage), input.SessionID, ref)

	core.CheckVersionDrift(cmd, sessionID)

	return nil
}
