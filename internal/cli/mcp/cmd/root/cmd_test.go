//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// TestMcpServe_FailsClosedOnUnsetCTXDIR is the regression guard
// required by spec/single-source-context-anchor.md. The MCP serve
// path must route through rc.RequireContextDir; with CTX_DIR
// unset, the cobra Run should return an error rather than starting
// a server bound to an empty path.
func TestMcpServe_FailsClosedOnUnsetCTXDIR(t *testing.T) {
	t.Setenv(env.CtxDir, "")
	rc.Reset()
	t.Cleanup(rc.Reset)

	c := &cobra.Command{Use: "serve"}
	c.SetArgs(nil)

	err := Cmd(c, nil)
	if err == nil {
		t.Fatal("Cmd() err = nil, want non-nil when CTX_DIR is unset")
	}
}
