//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"github.com/spf13/cobra"
)

var (
	watchLog    string
	watchDryRun bool
)

// Cmd returns the watch command.
//
// Flags:
//   - --log: Log file to watch (default: stdin)
//   - --dry-run: Show updates without applying
//
// Returns:
//   - *cobra.Command: Configured watch command with flags registered
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch for context-update commands in AI output",
		Long: `Watch stdin or a log file for <context-update>
commands and apply them.

This command parses AI output looking for structured update commands:

  Simple formats (tasks, conventions, complete):
    <context-update type="task">Implement user auth</context-update>
    <context-update type="convention">Use kebab-case for files</context-update>
    <context-update type="complete">user auth</context-update>

  Structured formats (learnings, decisions) - all attributes required:
    <context-update type="learning" context="Debugging hooks"
      lesson="Hooks receive JSON via stdin"
      application="Use jq to parse">Title here</context-update>

    <context-update type="decision" context="Need caching"
      rationale="Redis is fast and well-supported"
      consequences="Team needs Redis training">Use Redis</context-update>

Learnings require: context, lesson, application attributes.
Decisions require: context, rationale, consequences attributes.
Updates missing required attributes will be rejected with an error.

Use --log to watch a specific file instead of stdin.
Use --dry-run to see what would be updated without making changes.

Press Ctrl+C to stop watching.`,
		RunE: runWatch,
	}

	cmd.Flags().StringVar(
		&watchLog, "log", "", "Log file to watch (default: stdin)",
	)
	cmd.Flags().BoolVar(
		&watchDryRun, "dry-run", false, "Show updates without applying",
	)

	return cmd
}
