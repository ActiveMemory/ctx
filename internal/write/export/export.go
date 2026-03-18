//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Summary prints what an export will (or would) do based on
// aggregate counters.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - newCount: number of new files to export.
//   - regenCount: number of existing files to regenerate.
//   - skipCount: number of existing files to skip.
//   - lockedCount: number of locked files to skip.
//   - dryRun: when true, uses "Would" instead of "Will".
func Summary(
	cmd *cobra.Command,
	newCount, regenCount, skipCount, lockedCount int,
	dryRun bool,
) {
	if cmd == nil {
		return
	}

	verb := "Will"
	if dryRun {
		verb = "Would"
	}
	var parts []string
	if newCount > 0 {
		parts = append(parts, fmt.Sprintf("export %d new", newCount))
	}
	if regenCount > 0 {
		parts = append(parts, fmt.Sprintf("regenerate %d existing", regenCount))
	}
	if skipCount > 0 {
		parts = append(parts, fmt.Sprintf("skip %d existing", skipCount))
	}
	if lockedCount > 0 {
		parts = append(parts, fmt.Sprintf("skip %d locked", lockedCount))
	}
	if len(parts) == 0 {
		cmd.Println("Nothing to export.")
		return
	}
	cmd.Println(fmt.Sprintf("%s %s.", verb, strings.Join(parts, ", ")))
}
