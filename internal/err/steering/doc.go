//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package steering defines the **typed error
// constructors** returned by [internal/steering] —
// frontmatter parse failures, sync target validation,
// path-boundary violations, and missing-tool errors.
//
// # Why Typed Errors
//
//   - **Stability** — error categories are part of
//     the public API.
//   - **Routing** — write-side packages map error
//     types to localized text via
//     [internal/assets/read/desc].
//   - **Wrapping** — constructors wrap the
//     underlying cause via `%w` so callers can
//     `errors.Is` against system errors when
//     needed.
//
// # Public Surface
//
// Constructors fall into three groups:
//
//   - **Parse / IO** — [Parse], [InvalidYAML],
//     [MissingClosingDelimiter],
//     [MissingOpeningDelimiter], [ReadFile],
//     [ReadDir], [WriteFile], [WriteSteeringFile],
//     [WriteInitFile].
//   - **Sync** — [SyncAll], [SyncName],
//     [UnsupportedTool], [NoTool],
//     [ResolveOutput], [ResolveRoot],
//     [OutputEscapesRoot], [ComputeRelPath],
//     [CreateDir], [FileExists].
//   - **Context** — [ContextDirMissing].
//
// # The Boundary Check
//
// [OutputEscapesRoot] is fired when a sync
// target's resolved absolute path would land
// outside the project root — a defensive check
// that prevents a malicious or buggy steering
// file from writing to arbitrary filesystem
// locations.
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
//
// # Related Packages
//
//   - [internal/steering]              — chief
//     producer.
//   - [internal/cli/steering]          — also
//     producer.
//   - [internal/write/steering]        — renders
//     these into user-facing text.
package steering
