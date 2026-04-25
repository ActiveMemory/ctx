//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package post_commit

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/drift"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	corePC "github.com/ActiveMemory/ctx/internal/cli/system/core/post_commit"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	ctxContext "github.com/ActiveMemory/ctx/internal/context/resolve"
	"github.com/ActiveMemory/ctx/internal/notify"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Run executes the post-commit hook logic.
//
// After a successful git commit (non-amend), nudges the agent
// to offer context capture (decision or learning) and to run
// lints/tests before pushing. Also checks for version drift
// and spec enforcement.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	input, sessionID, _, _, ok := coreCheck.FullPreamble(stdin)
	bailSilently := !ok
	if bailSilently {
		return nil
	}

	command := input.ToolInput.Command

	if !regex.GitCommit.MatchString(command) {
		return nil
	}

	if regex.GitAmend.MatchString(command) {
		return nil
	}

	hookName := hook.PostCommit
	variant := hook.VariantNudge

	fallback := desc.Text(text.DescKeyPostCommitFallback)
	msg := message.Load(hookName, variant, nil, fallback)
	if msg == "" {
		return nil
	}
	msg, appendErr := ctxContext.AppendDir(msg)
	if appendErr != nil {
		return appendErr
	}
	writeSetup.Context(
		cmd,
		coreSession.FormatContext(hook.EventPostToolUse, msg),
	)

	ref := notify.NewTemplateRef(hookName, variant, nil)
	if relayErr := nudge.Relay(
		fmt.Sprintf(
			desc.Text(text.DescKeyRelayPrefixFormat),
			hookName,
			desc.Text(text.DescKeyPostCommitRelayMessage),
		),
		input.SessionID, ref,
	); relayErr != nil {
		return relayErr
	}

	driftResponse, driftErr := drift.CheckVersion(sessionID)
	if driftErr != nil {
		return driftErr
	}
	if driftResponse != "" {
		writeSetup.Context(cmd, driftResponse)
	}

	if violations := corePC.ScoreCommitViolations(); violations != "" {
		writeSetup.NudgeBlock(cmd, violations)
	}

	return nil
}
