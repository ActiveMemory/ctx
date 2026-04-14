//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package store manages **scratchpad file persistence**:
// the encrypted on-disk file, the AES-256-GCM key, and
// the read/write helpers every `ctx pad` subcommand uses.
//
// # Public Surface
//
//   - **[ScratchpadPath](contextDir)** — returns the
//     absolute path of `.context/.scratchpad.enc`.
//   - **[KeyPath]** — returns the per-machine key
//     path (`~/.ctx/.ctx.key`); shared with
//     [internal/notify].
//   - **[EnsureKey]** — creates the key on first use
//     and returns it. Subsequent calls just return
//     the existing key.
//   - **[ReadEntries](contextDir)** — decrypts the
//     scratchpad, parses it via
//     [internal/cli/pad/core/parse], returns a
//     `[]Entry`. Returns an empty slice (not an
//     error) when the file does not exist yet.
//   - **[WriteEntries](contextDir, entries)** —
//     formats, encrypts, atomically writes. Backup
//     is unnecessary because re-encryption never
//     produces a partial file when the rename
//     succeeds.
//
// # Atomic Writes
//
// Writes go through `tmpfile + os.Rename` so a
// crashed process never leaves a half-written
// `.scratchpad.enc`. The temp file is always created
// in the same directory as the destination so the
// rename is on the same filesystem (atomic on
// POSIX).
//
// # Concurrency
//
// Single-process. Concurrent writers would race on
// the temp filename collision; ctx is single-process
// by design.
//
// # Related Packages
//
//   - [internal/cli/pad/core/parse]   — parser used
//     during reads.
//   - [internal/crypto]               — encrypt /
//     decrypt primitives.
//   - [internal/notify]               — shares the
//     same per-machine key.
//   - [internal/cli/pad]              — top-level
//     CLI.
package store
