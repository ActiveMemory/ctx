//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_backup_age implements the
// **`ctx system check-backup-age`** hidden hook, which
// warns when project backups are stale or the backup
// share is unreachable.
//
// # What It Does
//
// The hook runs two checks at session start:
//
//  1. **SMB mount check**: if the CTX_BACKUP_SMB_URL
//     environment variable is set, verifies that the
//     SMB share is currently mounted. If not, adds a
//     warning.
//  2. **Backup marker freshness**: reads the backup
//     marker file (~/.ctx-backup/last-backup) and warns
//     if the last backup timestamp exceeds the maximum
//     age threshold.
//
// When either check produces warnings, the hook emits
// a nudge box via the relay channel and touches a
// daily throttle file to avoid repeated alerts.
//
// # Input
//
// A JSON hook envelope on stdin with session metadata.
//
// # Output
//
// On warning: a formatted nudge box listing the backup
// issues. On success or throttled: no output.
//
// # Throttling
//
// The hook is throttled to fire at most once per day
// using a marker file in the state directory.
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin via [core/check.Preamble], delegates the SMB
// and marker checks to [core/archive], and emits
// warnings through [core/nudge.LoadAndEmit].
package check_backup_age
