//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resume

import (
	"github.com/spf13/cobra"

	resumeroot "github.com/ActiveMemory/ctx/internal/cli/resume/cmd/root"
)

// Cmd returns the top-level "ctx resume" command.
func Cmd() *cobra.Command {
	return resumeroot.Cmd()
}
