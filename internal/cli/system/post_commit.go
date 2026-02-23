//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/notify"
)

// postCommitCmd returns the "ctx system post-commit" command.
//
// Fires after a successful git commit (PostToolUse on Bash). Detects git
// commit commands and nudges the agent to offer context capture and suggest
// running lints/tests before the user pushes.
func postCommitCmd() *cobra.Command {
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
	if !isInitialized() {
		return nil
	}
	input := readInput(stdin)
	command := input.ToolInput.Command

	// Only trigger on git commit commands
	if !reGitCommit.MatchString(command) {
		return nil
	}

	// Skip amend commits
	if reAmend.MatchString(command) {
		return nil
	}

	msg := "Commit succeeded." +
		" 1. Offer context capture to the user:" +
		" Decision (design choice?), Learning (gotcha?), or Neither." +
		" 2. Ask the user: \"Want me to run lints and tests before you push?\"" +
		" Do NOT push. The user pushes manually."
	if line := contextDirLine(); line != "" {
		msg += " [" + line + "]"
	}
	printHookContext(cmd, "PostToolUse", msg)

	_ = notify.Send("relay", "post-commit: Commit succeeded, context capture offered", input.SessionID)

	return nil
}
