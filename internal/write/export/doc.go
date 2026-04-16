//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package export provides terminal output for file
// export operations across multiple commands (recall
// import, pad export, and similar workflows).
//
// # Exported Functions
//
// [InfoExistsWritingAsAlternative] notifies the user
// when an output file already exists at the target
// path and an alternative filename is used to avoid
// overwriting. The message includes both the original
// path and the fallback path so the user knows where
// to find the output.
//
// # Message Categories
//
//   - Info: alternative filename notice
//
// # Nil Safety
//
// A nil *cobra.Command is treated as a no-op.
//
// # Usage
//
//	if fileExists(target) {
//	    alt := timestampedName(target)
//	    export.InfoExistsWritingAsAlternative(
//	        cmd, target, alt)
//	}
package export
