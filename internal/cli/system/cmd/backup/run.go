//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backup

import (
	"encoding/json"
	"os"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/env"
	backup2 "github.com/ActiveMemory/ctx/internal/err/backup"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/write/backup"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
)

// Run executes the backup command logic.
//
// Creates timestamped tar.gz archives of project context and/or global
// Claude Code data. Optionally copies archives to an SMB share.
//
// Parameters:
//   - cmd: Cobra command for output and flag access
//
// Returns:
//   - error: Non-nil on invalid scope, home directory lookup failure,
//     SMB parse error, or archive creation failure
func Run(cmd *cobra.Command) error {
	scope, _ := cmd.Flags().GetString("scope")
	jsonOut, _ := cmd.Flags().GetBool("json")

	switch scope {
	case archive.BackupScopeProject, archive.BackupScopeGlobal, archive.BackupScopeAll:
	default:
		return backup2.InvalidBackupScope(scope)
	}

	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return ctxerr.HomeDir(homeErr)
	}

	smbURL := os.Getenv(env.BackupSMBURL)
	smbSubdir := os.Getenv(env.BackupSMBSubdir)
	var smb *core.SMBConfig
	if smbURL != "" {
		var smbErr error
		smb, smbErr = core.ParseSMBConfig(smbURL, smbSubdir)
		if smbErr != nil {
			return backup2.SMBConfig(smbErr)
		}
	}

	timestamp := time.Now().Format(archive.BackupTimestampFormat)
	var results []core.BackupResult

	if scope == archive.BackupScopeProject || scope == archive.BackupScopeAll {
		result, projErr := core.BackupProject(cmd, home, timestamp, smb)
		if projErr != nil {
			return backup2.Project(projErr)
		}
		results = append(results, result)
	}

	if scope == archive.BackupScopeGlobal || scope == archive.BackupScopeAll {
		result, globalErr := core.BackupGlobal(cmd, home, timestamp, smb)
		if globalErr != nil {
			return backup2.Global(globalErr)
		}
		results = append(results, result)
	}

	if jsonOut {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}

	for _, r := range results {
		backup.ResultLine(cmd, r.Scope, r.Archive, r.Size, r.SMBDest)
	}
	return nil
}
