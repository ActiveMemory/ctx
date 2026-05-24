//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/audit/core/store"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeAudit "github.com/ActiveMemory/ctx/internal/write/audit"
)

// Run prints one audit report's body verbatim (no
// frontmatter, no formatting). Designed for unix-style
// piping into other tools.
//
// Parameters:
//   - cmd: Cobra command for output
//   - id: report basename (e.g. "surface")
//
// Returns:
//   - error: [errAudit.UnknownID] when the report does not
//     exist, or any underlying read / parse error
func Run(cmd *cobra.Command, id string) error {
	if _, ctxErr := rc.RequireContextDir(); ctxErr != nil {
		cmd.SilenceUsage = true
		return ctxErr
	}
	r, readErr := store.ReadOne(id)
	if readErr != nil {
		return readErr
	}
	writeAudit.Body(cmd, r.Body)
	return nil
}
