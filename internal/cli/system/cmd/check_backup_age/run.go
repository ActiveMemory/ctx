//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_backup_age

import (
	"os"
	"path/filepath"

	archive2 "github.com/ActiveMemory/ctx/internal/cli/system/core/archive"
	hook2 "github.com/ActiveMemory/ctx/internal/cli/system/core/hook"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// Run executes the check-backup-age hook logic.
//
// Reads a hook input from stdin, checks whether the SMB share is mounted
// and whether the backup marker file is fresh, then emits a relay warning
// if any issue is detected. Throttled to once per day.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	input, _, paused := hook2.Preamble(stdin)
	if paused {
		return nil
	}

	tmpDir := core.StateDir()
	throttleFile := filepath.Join(tmpDir, archive.BackupThrottleID)

	if core.DailyThrottled(throttleFile) {
		return nil
	}

	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return nil
	}

	var warnings []string

	// Check 1: Is the SMB share mounted?
	if smbURL := os.Getenv(env.BackupSMBURL); smbURL != "" {
		warnings = archive2.CheckSMBMountWarnings(smbURL, warnings)
	}

	// Check 2: Is the backup stale?
	markerPath := filepath.Join(
		home, archive.BackupMarkerDir, archive.BackupMarkerFile,
	)
	warnings = archive2.CheckBackupMarker(markerPath, warnings)

	if len(warnings) == 0 {
		return nil
	}

	// Build pre-formatted warnings for the template variable
	var warningText string
	for _, w := range warnings {
		warningText += w + token.NewlineLF
	}

	vars := map[string]any{archive.VarWarnings: warningText}
	content := core.LoadMessage(
		hook.CheckBackupAge, hook.VariantWarning, vars, warningText,
	)
	if content == "" {
		return nil
	}

	core.EmitNudge(cmd, content,
		desc.Text(text.DescKeyBackupRelayPrefix),
		desc.Text(text.DescKeyBackupBoxTitle),
		hook.CheckBackupAge, hook.VariantWarning,
		desc.Text(text.DescKeyBackupRelayMessage),
		input.SessionID, vars, throttleFile)

	return nil
}
