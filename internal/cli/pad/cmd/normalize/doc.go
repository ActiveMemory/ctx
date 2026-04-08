//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package normalize provides the pad normalize subcommand.
//
// It reassigns entry IDs as 1..N in current file order,
// closing all gaps left by deletions. This invalidates
// previously-seen IDs.
package normalize
