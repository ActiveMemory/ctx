//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package markwrappedup

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
)

// Cmd returns the "ctx system mark-wrapped-up" subcommand.
//
// Returns:
//   - *cobra.Command: Configured mark-wrapped-up subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-wrapped-up",
		Short: "Suppress checkpoint nudges after wrap-up",
		Long: `Write a marker file that suppresses context checkpoint nudges
for 2 hours. Called by /ctx-wrap-up after persisting context.

The check-context-size hook checks this marker before emitting
a checkpoint. If the marker exists and is less than 2 hours old,
the nudge is suppressed.

This is a plumbing command — use /ctx-wrap-up instead.`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runMarkWrappedUp(cmd)
		},
	}
}

// runMarkWrappedUp creates or updates the wrap-up marker file.
func runMarkWrappedUp(cmd *cobra.Command) error {
	markerPath := filepath.Join(core.StateDir(), core.WrappedUpMarker)

	if writeErr := os.WriteFile(
		markerPath, []byte("wrapped-up"), 0o600,
	); writeErr != nil {
		return writeErr
	}

	cmd.Println("marked wrapped-up")
	return nil
}
