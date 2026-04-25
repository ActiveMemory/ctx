//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/activate/core/emit"
	writeActivate "github.com/ActiveMemory/ctx/internal/write/activate"
)

// Run executes the `ctx deactivate` command: emit a shell-specific
// `unset CTX_DIR` statement to stdout so the caller can clear the
// binding via `eval "$(ctx deactivate)"`.
//
// The command never errors under normal operation; unsetting an
// already-unset variable is a no-op across supported shells.
//
// Parameters:
//   - cmd: cobra command providing stdout.
//   - shell: value of the --shell flag; empty auto-detects from
//     $SHELL via emit.DetectShell.
//
// Returns:
//   - error: always nil; kept in the signature for Cobra RunE
//     compatibility.
func Run(cmd *cobra.Command, shell string) error {
	writeActivate.Emit(cmd, emit.Unset(emit.DetectShell(shell)))
	return nil
}
