//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"github.com/spf13/cobra"

	statusRoot "github.com/ActiveMemory/ctx/internal/cli/status/cmd/root"
)

// Cmd returns the status command.
func Cmd() *cobra.Command {
	return statusRoot.Cmd()
}
