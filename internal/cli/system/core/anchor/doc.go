//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package anchor holds path-comparison helpers for the
// `ctx system check-anchor-drift` hook.
//
// The drift hook compares the parent-shell CTX_DIR (snapshotted
// before the standard `${CLAUDE_PROJECT_DIR:?…}/.context` injection)
// against the Claude-injected CTX_DIR. Naive byte-for-byte
// comparison after [filepath.Clean] over-reports: on macOS, `/tmp`
// is a symlink to `/private/tmp`, so a shell activated under
// `/tmp/foo/.context` and a Claude session whose
// `CLAUDE_PROJECT_DIR` resolves to `/private/tmp/foo` would trip
// a false alarm on every prompt: same physical directory, different
// strings. Any user with a symlinked workspace path runs into the
// same trap.
//
// # Public Surface
//
//   - [Equal] reports whether two paths refer to the same directory,
//     resolving symlinks before comparison and falling back to
//     cleaned-string comparison when resolution fails.
package anchor
