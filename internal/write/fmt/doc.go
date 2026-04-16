//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package fmt provides terminal output for the context
// file formatter command (ctx fmt).
//
// # Exported Functions
//
// [Summary] prints the formatting result showing how
// many context files were reformatted out of the total
// number scanned. This is printed to stdout after both
// format and check modes complete.
//
// [NeedsFormatting] prints a per-file message to stderr
// in check mode, identifying each context file that
// would need reformatting. This enables CI pipelines
// to detect formatting drift without modifying files.
//
// # Message Categories
//
//   - Info: format summary (stdout)
//   - Warning: per-file check-mode notices (stderr)
//
// # Usage
//
//	for _, f := range dirty {
//	    fmt.NeedsFormatting(cmd, f.Name)
//	}
//	fmt.Summary(cmd, len(dirty), len(all))
package fmt
