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
)

// postCommitCmd returns the "ctx system post-commit" command.
//
// Fires after a successful git commit (PostToolUse on Bash). Detects git
// commit commands and nudges the agent to offer context capture and suggest
// running lints/tests before the user pushes.
func postCommitCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "post-commit",
		Short:  "Post-commit context capture nudge",
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

	cmd.Println()
	cmd.Println("┌─ Post-Commit ──────────────────────────────────────────")
	cmd.Println("│ Commit succeeded.")
	cmd.Println("│")
	cmd.Println("│ 1. Offer context capture to the user:")
	cmd.Println("│    Decision (design choice?), Learning (gotcha?), or Neither.")
	cmd.Println("│")
	cmd.Println("│ 2. Ask the user:")
	cmd.Println("│    \"Want me to run lints and tests before you push?\"")
	cmd.Println("│")
	cmd.Println("│ Do NOT push. The user pushes manually.")
	cmd.Println("└────────────────────────────────────────────────────────")
	cmd.Println()

	return nil
}
