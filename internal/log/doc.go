//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package log provides event logging and stderr
// warning subpackages for ctx.
//
// This package itself contains no exported symbols;
// all functionality lives in its subpackages.
//
// # Subpackages
//
// [warn] provides a centralized stderr sink for
// best-effort operations whose errors would otherwise
// be silently discarded. Every non-fatal error in ctx
// flows through [warn.Warn] to keep warning output
// consistent across the codebase.
//
// The event subpackage writes and queries timestamped
// JSONL event logs for hook lifecycle tracking with
// automatic rotation. Events are written to
// .context/state/events.jsonl when event logging is
// enabled in .ctxrc.
//
// # Design Rationale
//
// ctx avoids the standard library's log package and
// third-party loggers. Instead, it uses structured
// JSONL for machine-readable events and a simple
// fprintf-to-stderr for human-readable warnings.
// This keeps the dependency surface small and gives
// each consumer explicit control over output format.
package log
