//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package validate checks that scratchpad entry
// indexes are within bounds.
//
// Several pad subcommands accept a 1-based entry index
// from the user (edit, delete, get). Before mutating
// or reading an entry, the index must be validated
// against the current entry count. This package
// provides that check.
//
// # Index Validation
//
// [Index] is the sole exported function. It accepts a
// 1-based index n and the current entries slice. If n
// is less than 1 or greater than len(entries), it
// returns an errPad.EntryRange error describing the
// valid range. Otherwise it returns nil.
//
// The edit, delete, and get functions in sibling
// packages call Index before accessing entries[n-1].
// This centralizes bounds checking so that each
// caller does not repeat the same logic.
//
// # Error Handling
//
// The returned error is a typed errPad.EntryRange
// value that includes the invalid index and the valid
// range, allowing the cmd layer to render a clear
// message to the user.
package validate
