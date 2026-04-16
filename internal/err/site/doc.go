//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package site defines the typed error constructors
// for the static site generation subsystem. These
// errors fire when the `ctx site` command validates
// the zensical configuration, marshals Atom feeds,
// or checks for the zensical binary.
//
// # Domain
//
// Three constructors cover the entire surface:
//
//   - [NoConfig]: the zensical.toml configuration
//     file is missing from the expected directory.
//   - [MarshalFeed]: the Atom XML feed could not
//     be marshaled from journal entries. Wraps the
//     underlying encoding/xml error.
//   - [ZensicalNotFound]: the zensical binary is
//     not installed. Returns a plain error with
//     installation instructions.
//
// # Wrapping Strategy
//
// [MarshalFeed] wraps its cause with fmt.Errorf %w.
// [NoConfig] returns a plain formatted error.
// [ZensicalNotFound] returns a plain errors.New
// value. All user-facing text is resolved through
// [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package site
