//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package change implements the ctx change command, which
// detects context and code changes since the last session
// or a specified time.
//
// The change command scans .context/ files and the git
// working tree for modifications, additions, and deletions
// that occurred after a reference timestamp. This helps
// AI agents and users understand what shifted between
// sessions without manually diffing files.
//
// # Detection Strategy
//
// The core/ subpackage walks context files and compares
// modification times against a cutoff. Git-tracked source
// files are checked via git diff. Results are grouped by
// change type (context vs code) and rendered as a
// human-readable summary or JSON.
//
// # Subpackages
//
//   - cmd/root: cobra command definition and flag binding
//   - core: change detection, scanning, and rendering
package change
