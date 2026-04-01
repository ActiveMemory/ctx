//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dep wraps dependency toolchain command execution.
//
// Centralizes exec.Command calls for Go and Rust toolchains.
// Callers get raw byte output and handle JSON decoding themselves.
//
// Key exports: [GoListPackages], [CargoMetadata].
// Part of the exec subsystem.
package dep
