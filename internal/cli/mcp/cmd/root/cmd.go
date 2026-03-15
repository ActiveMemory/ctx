//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	internalmcp "github.com/ActiveMemory/ctx/internal/mcp/server"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/spf13/cobra"
)

func Cmd(cmd *cobra.Command, _ []string) error {
	srv := internalmcp.NewServer(rc.ContextDir(), cmd.Root().Version)
	return srv.Serve()
}
