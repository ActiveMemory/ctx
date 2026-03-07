//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkbackupage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/notify"
)

const (
	backupMaxAgeDays = 2
	backupThrottleID = "backup-reminded"
)

// Cmd returns the "ctx system check-backup-age" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-backup-age subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-backup-age",
		Short: "Backup staleness check hook",
		Long: `Checks if the .context backup is stale (>2 days old) or the SMB share
is unmounted. Outputs a VERBATIM relay warning when issues are found.
Throttled to once per day.

Environment:
  CTX_BACKUP_SMB_URL - SMB share URL (e.g. smb://myhost/myshare).
                       If unset, the SMB mount check is skipped.

Hook event: UserPromptSubmit
Output: VERBATIM relay with warning box, silent otherwise
Silent when: backup is fresh, or already checked today`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckBackupAge(cmd, os.Stdin)
		},
	}
}

func runCheckBackupAge(cmd *cobra.Command, stdin *os.File) error {
	input := core.ReadInput(stdin)

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = core.SessionUnknown
	}
	if core.Paused(sessionID) > 0 {
		return nil
	}

	tmpDir := core.StateDir()
	throttleFile := filepath.Join(tmpDir, backupThrottleID)

	if core.IsDailyThrottled(throttleFile) {
		return nil
	}

	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return nil
	}

	var warnings []string

	// Check 1: Is the SMB share mounted?
	if smbURL := os.Getenv(config.EnvBackupSMBURL); smbURL != "" {
		warnings = checkSMBMountWarnings(smbURL, warnings)
	}

	// Check 2: Is the backup stale?
	markerPath := filepath.Join(home, ".local", "state", config.BackupMarkerFile)
	warnings = checkBackupMarker(markerPath, warnings)

	if len(warnings) == 0 {
		return nil
	}

	// Build pre-formatted warnings for the template variable
	var warningText string
	for _, w := range warnings {
		warningText += w + config.NewlineLF
	}

	content := core.LoadMessage("check-backup-age", "warning",
		map[string]any{"Warnings": warningText}, warningText)
	if content == "" {
		return nil
	}

	// Emit VERBATIM relay
	msg := "IMPORTANT: Relay this backup warning to the user VERBATIM before answering their question.\n\n" +
		"┌─ Backup Warning ──────────────────────────────────\n"
	msg += core.BoxLines(content)
	if line := core.ContextDirLine(); line != "" {
		msg += "│ " + line + config.NewlineLF
	}
	msg += config.NudgeBoxBottom
	cmd.Println(msg)

	ref := notify.NewTemplateRef("check-backup-age", "warning",
		map[string]any{"Warnings": warningText})
	_ = notify.Send("nudge", "check-backup-age: Backup warning", input.SessionID, ref)
	_ = notify.Send("relay", "check-backup-age: Backup warning", input.SessionID, ref)

	core.TouchFile(throttleFile)

	return nil
}

// checkSMBMountWarnings checks if the GVFS mount for the given SMB URL exists.
func checkSMBMountWarnings(smbURL string, warnings []string) []string {
	cfg, cfgErr := core.ParseSMBConfig(smbURL, "")
	if cfgErr != nil {
		return warnings
	}

	if _, statErr := os.Stat(cfg.GVFSPath); os.IsNotExist(statErr) {
		warnings = append(warnings,
			fmt.Sprintf("SMB share (%s) is not mounted.", cfg.Host),
			"Backups cannot run until it's available.",
		)
	}

	return warnings
}

// checkBackupMarker checks the backup marker file age.
func checkBackupMarker(markerPath string, warnings []string) []string {
	info, statErr := os.Stat(markerPath)
	if os.IsNotExist(statErr) {
		return append(warnings,
			"No backup marker found — backup may have never run.",
			"Run: ctx system backup",
		)
	}
	if statErr != nil {
		return warnings
	}

	ageDays := int(time.Since(info.ModTime()).Hours() / 24)
	if ageDays >= backupMaxAgeDays {
		return append(warnings,
			fmt.Sprintf("Last .context backup is %d days old.", ageDays),
			"Run: ctx system backup",
		)
	}

	return warnings
}
