//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package fs defines the **typed error constructors**
// for filesystem-level operations every other ctx
// package eventually performs: directory creation,
// reading, writing, amending. The package is the
// lowest level of the typed-error layer.
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
//     `errors.Is(err, os.ErrNotExist)` against
//     system errors when needed.
//
// # Public Surface
//
// Constructors: [Mkdir], [ReadDir], [DirNotFound],
// [FileWrite], [FileRead], [FileAmend].
//
// # When to Use [DirNotFound] vs [ReadDir]
//
// [DirNotFound] is for the actionable case "the
// directory the user expects to exist does not";
// [ReadDir] wraps the underlying generic read
// failure (permission denied, IO error, etc.).
// The CLI surfaces them differently — the former
// suggests `ctx init`, the latter suggests
// checking permissions.
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
//
// # Related Packages
//
//   - [internal/io]              — chief producer
//     for low-level filesystem failures.
//   - Most CLI commands           — surface them
//     when an expected file is missing.
package fs
