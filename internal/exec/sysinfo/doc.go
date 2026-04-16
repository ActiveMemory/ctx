//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sysinfo provides helpers for executing
// system information commands used by the sysinfo
// collector.
//
// # macOS Commands
//
// Sysctl runs the sysctl command with the given
// arguments and returns raw stdout output. This is
// used to query hardware parameters like memory
// size and CPU core count.
//
//	out, err := sysinfo.Sysctl("-n", "hw.memsize")
//
// VMStat runs the vm_stat command and returns raw
// stdout output. This is used to query virtual
// memory statistics like page counts and swap usage.
//
//	out, err := sysinfo.VMStat()
//
// # Build Constraints
//
// The implementation is gated behind a darwin build
// tag. Other platforms would need their own files
// with equivalent queries (e.g., reading /proc on
// Linux).
//
// # Centralization
//
// This package centralizes os/exec calls for
// platform-specific system queries, keeping nolint
// annotations in one place. The commands are fixed
// strings with no user input, but are routed through
// internal/exec to satisfy the project convention
// of no exec.Command calls outside this tree.
package sysinfo
