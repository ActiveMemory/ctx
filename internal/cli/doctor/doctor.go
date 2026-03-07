//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package doctor provides the "ctx doctor" command for structural
// health checks across context, hooks, and configuration.
package doctor

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/doctor/cmd/root"
	"github.com/ActiveMemory/ctx/internal/config"
)

// Cmd returns the "ctx doctor" command.
//
// Flags:
//   - --json, -j: Machine-readable JSON output
//
// Returns:
//   - *cobra.Command: Configured doctor command with flags registered
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "doctor",
		Short:       "Structural health check",
		Annotations: map[string]string{config.AnnotationSkipInit: "true"},
		Long: `Run mechanical health checks across context, hooks, and configuration.

Checks:
  - Context initialized and required files present
  - .ctxrc validation (unknown fields, typos)
  - Drift detected (stale paths, missing files)
  - Plugin installed and enabled
  - Event logging status
  - Webhook configured
  - Pending reminders
  - Task completion ratio
  - Context token size
  - System resources (memory, swap, disk, load)

Use --json for machine-readable output.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			jsonOut, _ := cmd.Flags().GetBool("json")
			return root.Run(cmd, jsonOut)
		},
	}
	cmd.Flags().BoolP("json", "j", false, "Machine-readable JSON output")
	return cmd
}
