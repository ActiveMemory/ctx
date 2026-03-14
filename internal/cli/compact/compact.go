//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"github.com/spf13/cobra"

	compactroot "github.com/ActiveMemory/ctx/internal/cli/compact/cmd/root"
)

// Cmd returns the "ctx compact" command for cleaning up context files.
func Cmd() *cobra.Command {
	return compactroot.Cmd()
}
