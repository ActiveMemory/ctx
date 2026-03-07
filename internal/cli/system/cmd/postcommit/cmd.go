//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package postcommit

import (
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// Cmd returns the "ctx system post-commit" subcommand.
//
// Returns:
//   - *cobra.Command: Configured post-commit subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "post-commit",
		Short: "Post-commit context capture nudge",
		Long: `Detects git commit commands and nudges the agent to offer context
capture (decision or learning) and suggest running lints/tests.
Skips amend commits.

Hook event: PostToolUse (Bash)
Output: agent directive after git commits, silent otherwise
Silent when: command is not a git commit, or is an amend`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPostCommit(cmd, os.Stdin)
		},
	}
}

var (
	reGitCommit = regexp.MustCompile(`git\s+commit`)
	reAmend     = regexp.MustCompile(`--amend`)
)

func runPostCommit(cmd *cobra.Command, stdin *os.File) error {
	if !core.IsInitialized() {
		return nil
	}
	input := core.ReadInput(stdin)

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = core.SessionUnknown
	}
	if core.Paused(sessionID) > 0 {
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

	fallback := "Commit succeeded." +
		" 1. Offer context capture to the user:" +
		" Decision (design choice?), Learning (gotcha?), or Neither." +
		" 2. Ask the user: \"Want me to run lints and tests before you push?\"" +
		" Do NOT push. The user pushes manually."
	msg := core.LoadMessage("post-commit", "nudge", nil, fallback)
	if msg == "" {
		return nil
	}
	if line := core.ContextDirLine(); line != "" {
		msg += " [" + line + "]"
	}
	core.PrintHookContext(cmd, "PostToolUse", msg)

	ref := notify.NewTemplateRef("post-commit", "nudge", nil)
	_ = notify.Send("relay", "post-commit: Commit succeeded, context capture offered", input.SessionID, ref)
	eventlog.Append("relay", "post-commit: Commit succeeded, context capture offered", input.SessionID, ref)

	core.CheckVersionDrift(cmd, sessionID)

	return nil
}
