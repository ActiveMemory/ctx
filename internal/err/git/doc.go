//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package git defines the typed error constructors
// for git repository detection and interaction.
// These errors fire when the CLI needs a git
// repository but either git is not installed or the
// working directory is not inside a repository.
//
// # Domain
//
// Two constructors cover the entire surface:
//
//   - [NotFound]: the git binary is not on PATH.
//     Returns a plain error with installation
//     guidance loaded from the assets catalog.
//   - [NotInRepo]: git rev-parse failed, meaning
//     the current directory is not inside a git
//     repository. Wraps the underlying exec error.
//
// # Wrapping Strategy
//
// [NotInRepo] wraps its cause with fmt.Errorf %w
// so callers can inspect the exec failure.
// [NotFound] returns a plain errors.New value.
// All user-facing text is resolved through
// [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package git
