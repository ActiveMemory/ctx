//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sysinfo provides helpers for executing system information
// commands (sysctl, vm_stat) used by the sysinfo collector.
//
// This package centralizes os/exec calls for platform-specific
// system queries, keeping nolint:gosec annotations in one place.
// The commands executed are fixed strings with no user input,
// but are routed through internal/exec/ to satisfy the project
// convention of no exec.Command calls outside this tree.
//
// Key exports: [Sysctl], [VMStat].
package sysinfo
