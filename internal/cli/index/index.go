//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package index

import (
	"github.com/spf13/cobra"

	indexRoot "github.com/ActiveMemory/ctx/internal/cli/index/cmd/root"
)

// Cmd returns the "ctx index" command, which projects the Markdown headings
// of a knowledge file as a computed table of contents.
//
// Returns:
//   - *cobra.Command: The index command.
func Cmd() *cobra.Command {
	return indexRoot.Cmd()
}
