//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package initialize hosts compile-time constants consumed by
// the ctx init command (backup directory naming, reset flag
// literal). Sentinel error values live in
// `internal/err/initialize/`; their user-facing text lives in
// `commands/text/errors.yaml` and is resolved through
// `desc.Text` at error-display time.
package initialize

// Backup directory naming for ctx init --reset.
const (
	// BackupDirPrefix is the basename prefix for timestamped
	// backup directories written by ctx init --reset before
	// it overwrites populated context files. The leading dot
	// keeps the directory out of glob-driven listings without
	// requiring .gitignore updates.
	BackupDirPrefix = ".backup-init-"
	// BackupTimestampLayout is the UTC ISO-8601 layout used
	// for backup directory suffixes (sortable and
	// filesystem-safe on every supported OS).
	BackupTimestampLayout = "2006-01-02T15-04-05Z"
	// BackupPlaceholder is the human-readable placeholder for
	// the backup directory when the actual timestamp is not
	// available (typically in error messages and docs).
	BackupPlaceholder = ".backup-init-<ISO>"
)

// ResetFlag is the canonical name of the destructive reset
// flag (--reset). Surfaces in error messages so prompts stay
// in sync with the cobra wiring even if the flag is renamed.
const ResetFlag = "--reset"
