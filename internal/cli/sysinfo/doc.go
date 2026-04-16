//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sysinfo implements the ctx sysinfo top-level
// command.
//
// Displays a snapshot of system resources including
// memory usage, swap utilization, disk space, and CPU
// load averages. Each metric is evaluated against
// configurable thresholds and assigned a severity level
// (ok, warning, danger) to surface resource pressure
// at a glance.
//
// # Output Formats
//
//   - Human-readable (default): a table with status
//     indicators (green/yellow/red) for each metric
//   - JSON (--json): a structured object containing
//     raw values and alert severities, suitable for
//     scripting and monitoring pipelines
//
// # How It Works
//
// [Run] delegates to [internal/cli/system/core/resource]
// for collecting the system snapshot and evaluating
// threshold alerts, then passes results to
// [internal/write/resource] for formatting.
//
// [Cmd] returns the cobra command with the --json flag.
// [Run] collects the system snapshot, evaluates each
// metric against its threshold, and writes the result
// in text or JSON format.
package sysinfo
