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

// Handler is the RunE installed on the `ctx system` group. It is
// reached only when cobra finds no matching subcommand (or for a
// bare `ctx system`); a valid subcommand runs its own RunE and never
// reaches here. It delegates to [handle] with the real stdin.
//
// Parameters:
//   - cmd: the system command (for output and SilenceUsage)
//   - args: leftover args; non-empty means an unknown subcommand
//
// Returns:
//   - error: nil for bare `ctx system` (help printed); otherwise the
//     unknown-subcommand error after emitting the relay box.
func Handler(cmd *cobra.Command, args []string) error {
	return handle(cmd, args, os.Stdin)
}
