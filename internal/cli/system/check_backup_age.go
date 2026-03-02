//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/notify"
)

const (
	backupMaxAgeDays = 2
	backupThrottleID = "backup-reminded"
)

// checkBackupAgeCmd returns the "ctx system check-backup-age" command.
//
// Warns when the SMB backup is stale (>2 days) or the share is unmounted.
// Depends on ctx system backup touching ~/.local/state/ctx-last-backup
// on successful backup. Throttled to once per day.
func checkBackupAgeCmd() *cobra.Command {
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
	input := readInput(stdin)

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = sessionUnknown
	}
	if paused(sessionID) > 0 {
		return nil
	}

	tmpDir := secureTempDir()
	throttleFile := filepath.Join(tmpDir, backupThrottleID)

	if isDailyThrottled(throttleFile) {
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
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
		warningText += w + "\n"
	}

	content := loadMessage("check-backup-age", "warning",
		map[string]any{"Warnings": warningText}, warningText)
	if content == "" {
		return nil
	}

	// Emit VERBATIM relay
	msg := "IMPORTANT: Relay this backup warning to the user VERBATIM before answering their question.\n\n" +
		"┌─ Backup Warning ──────────────────────────────────\n"
	msg += boxLines(content)
	if line := contextDirLine(); line != "" {
		msg += "│ " + line + "\n"
	}
	msg += config.NudgeBoxBottom
	cmd.Println(msg)

	ref := notify.NewTemplateRef("check-backup-age", "warning",
		map[string]any{"Warnings": warningText})
	_ = notify.Send("nudge", "check-backup-age: Backup warning", input.SessionID, ref)
	_ = notify.Send("relay", "check-backup-age: Backup warning", input.SessionID, ref)

	touchFile(throttleFile)

	return nil
}

// checkSMBMountWarnings checks if the GVFS mount for the given SMB URL exists.
func checkSMBMountWarnings(smbURL string, warnings []string) []string {
	cfg, cfgErr := parseSMBConfig(smbURL, "")
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
	info, err := os.Stat(markerPath)
	if os.IsNotExist(err) {
		return append(warnings,
			"No backup marker found — backup may have never run.",
			"Run: ctx system backup",
		)
	}
	if err != nil {
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
