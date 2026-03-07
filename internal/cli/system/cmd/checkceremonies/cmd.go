//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkceremonies

import (
	"os"

	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system check-ceremonies" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-ceremonies subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-ceremonies",
		Short: "Session ceremony nudge hook",
		Long: `Scans the last 3 journal entries for /ctx-remember and /ctx-wrap-up
usage. If either is missing, emits a VERBATIM relay nudge encouraging
adoption. Throttled to once per day.

Hook event: UserPromptSubmit
Output: VERBATIM relay (when ceremonies missing), silent otherwise
Silent when: both ceremonies found in recent sessions`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckCeremonies(cmd, os.Stdin)
		},
	}
}
