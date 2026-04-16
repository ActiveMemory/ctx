//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resource collects system resource metrics and
// evaluates alert thresholds. It provides a single entry
// point for hooks that need to check system health before
// or during a session.
//
// # Snapshot Collection
//
// [Snapshot] delegates to the sysinfo package to collect
// current system metrics (CPU, memory, disk) and then
// evaluates each metric against configured thresholds.
// It returns both the raw snapshot and a list of alerts
// for any metrics that exceed their thresholds.
//
// The two-step design separates collection from policy:
// sysinfo.Collect gathers raw data, sysinfo.Evaluate
// applies threshold rules. This package composes them
// into a single call for hook convenience.
//
// # Usage Pattern
//
// Hooks call Snapshot once per prompt cycle. When alerts
// are returned, the hook formats them into a nudge box
// warning the agent about resource pressure (e.g., low
// disk space or high memory usage).
package resource
