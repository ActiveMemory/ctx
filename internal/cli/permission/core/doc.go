//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides business logic for the permission
// command, which audits and sanitizes Claude Code
// settings files.
//
// The core package contains one subpackage, diff, which
// computes set differences between string slices. The
// permission command uses this to compare golden
// (reference) settings against local settings and
// identify entries that were restored or dropped.
//
// # Data Flow
//
// The cmd/ layer reads the golden and local settings
// files, then calls into core/diff to compute the
// difference. The results are passed to the write/
// layer for formatted output.
//
// # Key Subpackage
//
// The diff subpackage exports [diff.StringSlices], which
// takes two string slices (golden and local) and returns
// two slices: restored entries (in golden but not local)
// and dropped entries (in local but not golden). Both
// outputs preserve the ordering of their respective
// inputs.
package core
