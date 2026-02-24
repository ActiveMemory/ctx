//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sysinfo gathers OS-level resource metrics (memory, swap, disk, load)
// and evaluates them against configurable thresholds to produce alerts at
// WARNING and DANGER severity levels.
//
// Platform support uses build tags: Linux reads /proc, macOS shells out to
// sysctl/vm_stat, and other platforms return Supported: false gracefully.
package sysinfo
