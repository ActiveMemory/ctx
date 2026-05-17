//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package row supplies the shared append-a-monotonic-ID-row
// primitive used by every kb tabular artifact whose rows carry
// a `<PREFIX>-###` ID: contradictions, domain-decisions, and
// outstanding-questions. Evidence and source-coverage have
// distinct enough surfaces to stay in their own packages.
//
// The primitive [Append] handles the invariants that are
// identical across these artifacts:
//
//   - Ensure the parent directory exists.
//   - Read the existing file; create-on-not-exist.
//   - Delegate next-ID allocation to a caller-supplied [Hooks]
//     closure that knows how to scan its own table.
//   - Open the file O_CREATE|O_APPEND|O_WRONLY.
//   - On first write, emit the table header.
//   - Render and append the row via a caller-supplied closure.
//
// Each caller package owns its own Row struct, header
// constant, ID-scan logic, and renderer; this package is the
// glue that keeps the I/O shape consistent.
package row
