//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkversion

import (
	"os"

	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system check-version" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-version subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-version",
		Short: "Binary/plugin version drift detection hook",
		Long: `Compares the ctx binary version against the embedded plugin version.
Warns when the binary is older than the plugin expects, which happens
when the marketplace plugin updates but the binary hasn't been
reinstalled. Throttled to once per day. Skipped for dev builds.

Hook event: UserPromptSubmit
Output: VERBATIM relay with reinstall command, silent otherwise
Silent when: versions match, dev build, or already checked today`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckVersion(cmd, os.Stdin)
		},
	}
}
