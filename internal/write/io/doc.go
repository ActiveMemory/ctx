//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package io provides low-level print primitives shared across write
// subpackages. It is not intended for direct use by callers outside
// internal/write/; domain write packages wrap these primitives with
// domain-specific function names.
//
// Example usage from a domain write package:
//
//	// write/events/events.go
//	func JSON(cmd *cobra.Command, lines []string) {
//	    writeIO.Lines(cmd, lines)
//	}
//
//	// write/stat/stat.go
//	func Table(cmd *cobra.Command, lines []string) {
//	    writeIO.Lines(cmd, lines)
//	}
package io
