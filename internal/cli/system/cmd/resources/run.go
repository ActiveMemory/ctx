//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resources

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/sysinfo"
)

// runResources executes the resources display logic.
//
// Collects a system resource snapshot, evaluates alerts, and outputs
// results as either a JSON object or a human-readable table with
// status indicators.
//
// Parameters:
//   - cmd: Cobra command for output and flag access
//
// Returns:
//   - error: Non-nil on JSON encoding failure
func runResources(cmd *cobra.Command) error {
	snap := sysinfo.Collect(".")
	alerts := sysinfo.Evaluate(snap)

	jsonFlag, _ := cmd.Flags().GetBool("json")
	if jsonFlag {
		return core.OutputResourcesJSON(cmd, snap, alerts)
	}

	core.OutputResourcesText(cmd, snap, alerts)
	return nil
}
