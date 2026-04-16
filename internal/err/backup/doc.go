//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package backup defines the **typed error constructors**
// the backup subsystem returns. Every failure that can
// happen during `ctx backup` flows through one of these
// constructors so the call site upstream sees a
// sentinel-able error and the renderer downstream knows
// which user-facing text to surface.
//
// # Why Typed Errors
//
// Three reasons:
//
//   - **Stability**: error categories are part of
//     the public API; adding a constructor is an
//     intentional change a reviewer can see.
//   - **Routing**: the write-side
//     ([internal/write/backup]) maps error types to
//     localized text via [internal/assets/read/desc].
//   - **Wrapping**: every constructor wraps its
//     underlying cause via `%w` so callers can
//     `errors.Is` / `errors.As` against system
//     errors when needed.
//
// # Public Surface
//
// Constructors (one per failure mode): [Create],
// [CreateArchive], [CreateArchiveDir],
// [WriteArchive], [SMBConfig].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package backup
