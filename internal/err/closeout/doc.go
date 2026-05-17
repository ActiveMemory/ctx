//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package closeout defines the typed error constructors for
// Phase KB closeout artifacts. Closeouts are per-pass audit
// records written by the editorial pipeline under
// `.context/ingest/closeouts/`; this package owns every
// error surface the closeout reader / writer / archiver
// can return.
//
// # Domain
//
// Three sentinels + nine wrapping constructors cover the
// writer's full failure surface:
//
//   - [ErrMissingFrontmatter]: file missing the opening
//     `---` delimiter.
//   - [ErrMissingFields]: frontmatter parsed but missing
//     sha / branch / mode / generated-at.
//   - [ErrModeRequired]: Write was called with an empty
//     mode.
//   - [ReadFailed], [ParseFrontmatter], [MarshalFrontmatter],
//     [ReadCloseoutsDir], [ResolveHead], [MkdirCloseouts],
//     [WriteCloseout], [MkdirArchive], [ArchiveMove] wrap
//     I/O, YAML, and git failures with operator-friendly
//     context.
//
// # Wrapping strategy
//
// The constructor functions use `fmt.Errorf` with `%w` so
// callers can `errors.Is` against the sentinel where
// applicable and `errors.Unwrap` to recover the underlying
// cause for diagnostic output.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/closeout]
//     supplies the sentinel message and format-string
//     constants.
//   - [github.com/ActiveMemory/ctx/internal/write/closeout]
//     is the primary caller.
package closeout
