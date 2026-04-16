//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package line provides low-level print primitives
// shared across write subpackages.
//
// This package is not intended for direct use by
// callers outside internal/write/. Domain write
// packages wrap these primitives with domain-specific
// function names to keep their public API descriptive.
//
// # Exported Functions
//
// [All] prints each line in a string slice to the
// command's output. A nil *cobra.Command is a no-op.
// This eliminates the repeated nil-guard + range loop
// pattern across write packages.
//
// [Count] prints a formatted count line only when the
// count is positive. It accepts a descriptor key and
// formats the message with the count value. This
// eliminates the repeated if-count-gt-zero-print
// pattern used for summary lines.
//
// # Message Categories
//
//   - Info: line-by-line output and conditional counts
//
// # Usage
//
//	// From write/events:
//	func JSON(cmd *cobra.Command, lines []string) {
//	    line.All(cmd, lines)
//	}
//
//	// From write/publish:
//	line.Count(cmd, text.DescKeyPublishTasks, tasks)
package line
