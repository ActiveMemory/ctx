//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

// Backup configuration.
const (
	// BackupDefaultSubdir is the default subdirectory on the SMB share.
	BackupDefaultSubdir = "ctx-sessions"
	// BackupMarkerFile is the state file touched on a successful project backup.
	BackupMarkerFile = "ctx-last-backup"
	// BackupScopeProject backs up only the project context.
	BackupScopeProject = "project"
	// BackupScopeGlobal backs up only global Claude data.
	BackupScopeGlobal = "global"
	// BackupScopeAll backs up both project and global.
	BackupScopeAll = "all"
	// TplProjectArchive is the filename template for project archives.
	// Argument: timestamp.
	TplProjectArchive = "ctx-backup-%s.tar.gz"
	// TplGlobalArchive is the filename template for global archives.
	// Argument: timestamp.
	TplGlobalArchive = "claude-global-backup-%s.tar.gz"
	// BackupTimestampFormat is the compact timestamp layout for backup filenames.
	BackupTimestampFormat = "20060102-150405"
	// BackupExcludeTodos is the directory name excluded from global backups.
	BackupExcludeTodos = "todos"
	// BackupMarkerDir is the XDG state directory for the backup marker.
	BackupMarkerDir = ".local/state"
	// BackupMaxAgeDays is the threshold in days before a backup is considered stale.
	BackupMaxAgeDays = 2
	// BackupThrottleID is the state file name for daily throttle of backup age checks.
	BackupThrottleID = "backup-reminded"
	// Bashrc is the user's bash configuration file.
	Bashrc = ".bashrc"
)
