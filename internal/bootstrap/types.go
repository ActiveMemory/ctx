//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import "github.com/spf13/cobra"

// registration pairs a command constructor with its group ID.
// An empty groupID marks the command as hidden.
type registration struct {
	cmd     func() *cobra.Command
	groupID string
}
