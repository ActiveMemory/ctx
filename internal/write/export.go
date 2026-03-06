//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// ExportCounts holds aggregate counters for export summary output.
type ExportCounts struct {
	New    int
	Regen  int
	Skip   int
	Locked int
}

// ExportSummary prints what an export will (or would) do based on
// aggregate counters.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - counts: aggregate export counters.
//   - dryRun: when true, uses "Would" instead of "Will".
func ExportSummary(cmd *cobra.Command, counts ExportCounts, dryRun bool) {
	if cmd == nil {
		return
	}

	verb := "Will"
	if dryRun {
		verb = "Would"
	}
	var parts []string
	if counts.New > 0 {
		parts = append(parts, fmt.Sprintf("export %d new", counts.New))
	}
	if counts.Regen > 0 {
		parts = append(parts, fmt.Sprintf("regenerate %d existing", counts.Regen))
	}
	if counts.Skip > 0 {
		parts = append(parts, fmt.Sprintf("skip %d existing", counts.Skip))
	}
	if counts.Locked > 0 {
		parts = append(parts, fmt.Sprintf("skip %d locked", counts.Locked))
	}
	if len(parts) == 0 {
		cmd.Println("Nothing to export.")
		return
	}
	sprintf(cmd, "%s %s.", verb, strings.Join(parts, ", "))
}
