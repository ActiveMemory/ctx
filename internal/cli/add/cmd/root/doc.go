//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the ctx add command.
//
// [Cmd] builds the cobra.Command with type-specific flags.
// [Run] validates arguments, extracts content from args or
// --from-file, formats the entry using core/format, and inserts
// it into the target context file using core/insert.
package root
