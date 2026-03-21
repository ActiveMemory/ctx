//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	internalMcp "github.com/ActiveMemory/ctx/internal/mcp/server"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func Cmd(cmd *cobra.Command, _ []string) error {
	srv := internalMcp.NewServer(rc.ContextDir(), cmd.Root().Version)
	return srv.Serve()
}
