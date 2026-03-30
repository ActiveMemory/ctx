//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package plan builds the import plan for journal import operations.
//
// [Import] scans available sessions, matches them against existing
// journal files, and produces an ImportPlan describing what to
// create, regenerate, skip, or rename. Locked entries are always
// preserved.
package plan
