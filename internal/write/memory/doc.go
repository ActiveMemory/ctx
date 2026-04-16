//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package memory provides terminal output for the
// memory bridge commands (ctx memory status, sync, diff).
//
// # Status Dashboard
//
// Functions are composable: the caller assembles the
// status display by calling them in sequence with
// [StatusSeparator] between sections.
//
// [BridgeHeader] prints the "Memory Bridge Status"
// heading. [Source] prints the MEMORY.md source path.
// [SourceNotActive] prints a notice when auto memory
// is not active. [Mirror] prints the mirror relative
// path. [LastSync] prints the last sync timestamp with
// a human-readable age string. [LastSyncNever] prints
// that no sync has occurred yet.
//
// [SourceLines] prints the MEMORY.md line count with
// an optional drift indicator. [MirrorLines] prints
// the mirror line count. [MirrorNotSynced] prints
// that the mirror has not been synced yet.
//
// # Drift Detection
//
// [DriftDetected] prints that drift was detected
// between source and mirror. [DriftNone] prints that
// no drift was detected.
//
// # Archive and Diff
//
// [Archives] prints the archive snapshot count and
// directory. [DiffOutput] prints diff content to
// stdout. [NoChanges] prints that no changes exist
// since the last sync.
//
// # Message Categories
//
//   - Info: dashboard sections, sync status, counts
//   - Warning: drift detection notices
package memory
