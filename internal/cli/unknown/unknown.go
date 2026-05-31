//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unknown

import (
	"os"

	"github.com/spf13/cobra"
)

// HandlerFor returns a cobra RunE that relays an unknown subcommand for
// the group described by cfg. A group opts in by assigning the result
// to its RunE; cobra reaches it only when no subcommand matches (or for
// a bare group invocation) — a valid subcommand runs its own RunE and
// never reaches here. The returned closure delegates to [handle] with
// the real stdin.
//
// Parameters:
//   - cfg: the group's text keys and relay ref (e.g. [SystemConfig],
//     [HookConfig])
//
// Returns:
//   - func(*cobra.Command, []string) error: a RunE bound to cfg
func HandlerFor(cfg Config) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return handle(cmd, args, os.Stdin, cfg)
	}
}
