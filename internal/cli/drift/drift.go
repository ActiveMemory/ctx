//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"github.com/spf13/cobra"

	driftRoot "github.com/ActiveMemory/ctx/internal/cli/drift/cmd/root"
)

// Cmd returns the "ctx drift" command for detecting stale context.
func Cmd() *cobra.Command {
	return driftRoot.Cmd()
}
