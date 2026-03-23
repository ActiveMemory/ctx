//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resources

import (
	"github.com/spf13/cobra"

	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/sysinfo"
	writeResources "github.com/ActiveMemory/ctx/internal/write/resource"
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

	jsonFlag, _ := cmd.Flags().GetBool(cFlag.JSON)
	if jsonFlag {
		return writeResources.JSON(cmd, snap, alerts)
	}

	writeResources.Text(cmd, snap, alerts)
	return nil
}
