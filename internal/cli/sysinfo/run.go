//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import (
	"github.com/spf13/cobra"

	coreResource "github.com/ActiveMemory/ctx/internal/cli/system/core/resource"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	writeResources "github.com/ActiveMemory/ctx/internal/write/resource"
)

// Run executes the resources display logic.
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
func Run(cmd *cobra.Command) error {
	snap, alerts := coreResource.Snapshot()

	jsonFlag, _ := cmd.Flags().GetBool(cFlag.JSON)
	if jsonFlag {
		return writeResources.JSON(cmd, snap, alerts)
	}

	writeResources.Text(cmd, snap, alerts)
	return nil
}
