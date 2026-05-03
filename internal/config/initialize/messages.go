//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package initialize hosts compile-time constants consumed by
// the ctx init command (sentinel error messages, backup
// directory naming, reset flag literal).
//
// Sentinel error messages live here — not in the embedded YAML
// loaded via desc.Text — because the err/initialize package
// instantiates them at package-load time, before the YAML
// lookup table is populated.
package initialize

// Sentinel error messages for ctx init refusal and reset.
//
// These mirror keys in commands/text/errors.yaml but exist as
// raw string constants because the var ErrContextPopulated /
// var ErrResetRequiresInteractive sentinels in the err package
// are initialized before the YAML lookup is ready.
const (
	// ErrMsgContextPopulated is the sentinel message for
	// ctx init's refuse-when-populated guard.
	ErrMsgContextPopulated = "context already populated; refusing to overwrite"
	// ErrMsgResetRequiresInteractive is the sentinel message
	// for ctx init --reset's interactive-only guard.
	ErrMsgResetRequiresInteractive = "ctx init --reset requires" +
		" an interactive terminal"
)

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
