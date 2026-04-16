//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package backup implements the ctx backup top-level command.
//
// Creates timestamped tar.gz archives of project context
// and/or global Claude Code data. The --scope flag selects
// what to archive: "project" (the .context/ directory),
// "global" (Claude Code config under ~/.claude), or "all"
// (both). Archives are written to a local directory and
// optionally copied to an SMB share configured via the
// CTX_BACKUP_SMB_URL environment variable.
//
// The command delegates archive creation to
// [internal/cli/system/core/archive] and writes results
// through [internal/write/backup]. JSON output is
// available via --json for scripting.
//
// # Flags
//
//   - --scope: project | global | all (default: all)
//   - --json: machine-readable output
//
// # SMB Remote Copy
//
// When CTX_BACKUP_SMB_URL is set, the archive is copied
// to the configured share after local creation. An
// optional CTX_BACKUP_SMB_SUBDIR narrows the target
// directory within the share.
//
// [Cmd] builds the cobra command with scope and JSON flags.
// [Run] creates the archive for the selected scope and copies
// it to the SMB share when configured.
package backup
