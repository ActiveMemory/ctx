//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	internalmcp "github.com/ActiveMemory/ctx/internal/mcp"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/spf13/cobra"
)

func Cmd(_ *cobra.Command, _ []string) error {
	srv := internalmcp.NewServer(rc.ContextDir())
	return srv.Serve()
}
