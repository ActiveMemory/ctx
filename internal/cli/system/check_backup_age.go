//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/notify"
)

const (
	backupMaxAgeDays = 2
	backupMarkerFile = "ctx-last-backup"
	backupThrottleID = "backup-reminded"
)

// checkBackupAgeCmd returns the "ctx system check-backup-age" command.
//
// Warns when the SMB backup is stale (>2 days) or the share is unmounted.
// Depends on hack/backup-context.sh touching ~/.local/state/ctx-last-backup
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
	if smbURL := os.Getenv("CTX_BACKUP_SMB_URL"); smbURL != "" {
		warnings = checkSMBMount(smbURL, warnings)
	}

	// Check 2: Is the backup stale?
	markerPath := filepath.Join(home, ".local", "state", backupMarkerFile)
	warnings = checkBackupMarker(markerPath, warnings)

	if len(warnings) == 0 {
		return nil
	}

	// Emit VERBATIM relay
	cmd.Println("IMPORTANT: Relay this backup warning to the user VERBATIM before answering their question.")
	cmd.Println()
	cmd.Println("┌─ Backup Warning ──────────────────────────────────")
	for _, w := range warnings {
		cmd.Println("│ " + w)
	}
	if line := contextDirLine(); line != "" {
		cmd.Println("│ " + line)
	}
	cmd.Println("└──────────────────────────────────────────────────")

	_ = notify.Send("nudge", "check-backup-age: Backup warning", input.SessionID, "")
	_ = notify.Send("relay", "check-backup-age: Backup warning", input.SessionID, "")

	touchFile(throttleFile)

	return nil
}

// checkSMBMount checks if the GVFS mount for the given SMB URL exists.
func checkSMBMount(smbURL string, warnings []string) []string {
	u, err := url.Parse(smbURL)
	if err != nil || u.Host == "" {
		return warnings
	}

	host := u.Host
	share := u.Path
	if len(share) > 0 && share[0] == '/' {
		share = share[1:]
	}
	if share == "" {
		return warnings
	}

	gvfsPath := fmt.Sprintf("/run/user/%d/gvfs/smb-share:server=%s,share=%s",
		os.Getuid(), host, share)

	if _, err := os.Stat(gvfsPath); os.IsNotExist(err) {
		warnings = append(warnings,
			fmt.Sprintf("SMB share (%s) is not mounted.", host),
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
			"Run: hack/backup-context.sh",
		)
	}
	if err != nil {
		return warnings
	}

	ageDays := int(time.Since(info.ModTime()).Hours() / 24)
	if ageDays >= backupMaxAgeDays {
		return append(warnings,
			fmt.Sprintf("Last .context backup is %d days old.", ageDays),
			"Run: hack/backup-context.sh",
		)
	}

	return warnings
}
