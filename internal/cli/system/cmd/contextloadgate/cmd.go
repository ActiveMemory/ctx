//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package contextloadgate

import (
	"os"

	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system context-load-gate" subcommand.
//
// Returns:
//   - *cobra.Command: Configured context-load-gate subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "context-load-gate",
		Short: "Auto-inject project context on first tool use",
		Long: `Auto-injects project context into the agent's context window.
Fires on the first tool use per session via PreToolUse hook. Subsequent
tool calls in the same session are silent (tracked by session marker file).

Reads context files directly and injects content — no delegation to
bootstrap command, no agent compliance required.
See specs/context-load-gate-v2.md for design rationale.

Hook event: PreToolUse (.*)
Output: JSON HookResponse (additionalContext) on first tool use, silent otherwise
Silent when: marker exists for session_id, or context not initialized`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runContextLoadGate(cmd, os.Stdin)
		},
	}
}
