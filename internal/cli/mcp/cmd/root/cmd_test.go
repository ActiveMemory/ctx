//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/rc"
)

// TestMcpServe_FailsClosedWhenNoDotContext is the regression guard
// for the cwd-anchored model (specs/cwd-anchored-context.md). The
// MCP serve path must route through rc.RequireContextDir; when
// $PWD/.context/ is absent, the cobra Run should return an error
// rather than starting a server bound to an empty path.
func TestMcpServe_FailsClosedWhenNoDotContext(t *testing.T) {
	t.Chdir(t.TempDir())
	rc.Reset()
	t.Cleanup(rc.Reset)

	c := &cobra.Command{Use: "serve"}
	c.SetArgs(nil)

	err := Cmd(c, nil)
	if err == nil {
		t.Fatal("Cmd() err = nil, want non-nil when $PWD has no .context/")
	}
}
