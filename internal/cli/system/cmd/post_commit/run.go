//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package post_commit

import (
	"fmt"
	"os"
	"regexp"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	ctxcontext "github.com/ActiveMemory/ctx/internal/context/resolve"
	"github.com/spf13/cobra"

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
	if !core.Initialized() {
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

	hookName, variant := hook.PostCommit, hook.VariantNudge

	fallback := desc.Text(text.DescKeyPostCommitFallback)
	msg := core.LoadMessage(hookName, variant, nil, fallback)
	if msg == "" {
		return nil
	}
	msg = ctxcontext.AppendDir(msg)
	core.PrintHookContext(cmd, hook.EventPostToolUse, msg)

	ref := notify.NewTemplateRef(hookName, variant, nil)
	core.Relay(fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat), hookName, desc.Text(text.DescKeyPostCommitRelayMessage)), input.SessionID, ref)

	core.CheckVersionDrift(cmd, sessionID)

	return nil
}
