//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core contains drift detection output and auto-fix logic.
//
// [OutputDriftText] renders a drift report as a human-readable
// checklist. [ApplyFixes] runs interactive fixes for detected
// issues. [FixStaleness] touches stale files and [FixMissingFile]
// creates missing context files from templates.
package core
