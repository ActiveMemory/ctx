//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

// CheckOpts holds the flags for the schema check command.
//
// Fields:
//   - Dir: Custom directory to scan for JSONL files
//   - AllProjects: Scan all Claude Code project directories
//   - Quiet: Exit code only, no output
type CheckOpts struct {
	Dir         string
	AllProjects bool
	Quiet       bool
}
