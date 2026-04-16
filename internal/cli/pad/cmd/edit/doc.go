//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package edit implements **`ctx pad edit`** — the
// subcommand that decrypts the scratchpad to a temp file,
// spawns the user's `$EDITOR` against it, and re-encrypts
// the result on save.
//
// # Behavior
//
//   - **Cleartext temp file** lives in the secure
//     temp directory and is `0o600`.
//   - **Editor invocation** uses `$EDITOR` (or
//     `vi` as fallback). Foreground; ctx blocks
//     until the editor exits.
//   - **Re-encrypt** on successful exit. Editor
//     non-zero exit aborts the write and leaves the
//     scratchpad untouched.
//   - **Cleanup** — the temp file is removed in a
//     deferred handler regardless of outcome so a
//     crashed editor does not leak plaintext.
//
// # Concurrency
//
// Single-process, sequential.
package edit
