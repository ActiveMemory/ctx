//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package log provides event logging and stderr warning subpackages.
//
// Subpackage [event] writes and queries timestamped JSONL event logs
// for hook lifecycle tracking with automatic rotation. Subpackage
// [warn] provides a centralized stderr sink for best-effort operations
// whose errors would otherwise be silently discarded.
//
// This package itself contains no exported symbols; all functionality
// lives in the subpackages.
package log
