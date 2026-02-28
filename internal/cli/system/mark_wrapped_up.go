//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// wrappedUpMarker is the filename for the wrap-up suppression marker.
const wrappedUpMarker = "ctx-wrapped-up"

// wrappedUpExpiry is how long the marker suppresses nudges.
const wrappedUpExpiry = 2 * time.Hour

// markWrappedUpCmd returns the "ctx system mark-wrapped-up" subcommand.
//
// Writes a marker file to secureTempDir() that suppresses context
// checkpoint nudges for 2 hours. Called by the /ctx-wrap-up skill
// after persisting session context, so the wrap-up ceremony itself
// does not trigger noisy checkpoint reminders.
//
// Hidden because it is a plumbing command called by skills, not a
// user-facing workflow.
func markWrappedUpCmd() *cobra.Command {
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
	markerPath := filepath.Join(secureTempDir(), wrappedUpMarker)

	if writeErr := os.WriteFile(
		markerPath, []byte("wrapped-up"), 0o600,
	); writeErr != nil {
		return writeErr
	}

	cmd.Println("marked wrapped-up")
	return nil
}

// wrappedUpRecently checks whether the wrap-up marker exists and is
// less than wrappedUpExpiry old.
//
// Returns true if nudges should be suppressed.
func wrappedUpRecently() bool {
	markerPath := filepath.Join(secureTempDir(), wrappedUpMarker)

	info, statErr := os.Stat(markerPath)
	if statErr != nil {
		return false
	}

	return time.Since(info.ModTime()) < wrappedUpExpiry
}
