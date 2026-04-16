//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_freshness implements the
// **`ctx system check-freshness`** hidden hook, which
// warns when tracked source files have not been
// reviewed within the configured freshness window.
//
// # What It Does
//
// The hook reads the freshness_files list from .ctxrc.
// For each configured file, it stats the file on disk
// and compares the modification time against the
// staleness threshold (approximately six months). Files
// that do not exist are silently skipped. When one or
// more files exceed the threshold, the hook emits a
// nudge listing each stale file with its age in days,
// description, and review URL.
//
// This is useful for files containing technology-
// dependent constants, dependency versions, or external
// API references that may drift over time.
//
// # Input
//
// A JSON hook envelope on stdin with session metadata.
//
// # Output
//
// On stale files: a formatted nudge box listing each
// stale file path, its age, and a review link.
// On no stale files, no configured files, or throttled:
// no output.
//
// # Throttling
//
// The hook is throttled to fire at most once per day
// using a marker file in the state directory.
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin via [core/check.Preamble], iterates over
// [rc.FreshnessFiles], formats stale entries through
// [core/drift.FormatStaleEntries], and emits the nudge
// via [core/nudge.LoadAndEmit].
package check_freshness
