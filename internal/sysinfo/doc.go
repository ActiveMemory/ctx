//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sysinfo gathers OS-level resource metrics — memory,
// swap, disk, load average — and evaluates them against
// configurable thresholds to produce alerts at **WARNING** and
// **DANGER** severity levels.
//
// The package powers two surfaces:
//
//   - **`ctx sysinfo`** — the top-level user-facing CLI that
//     prints a snapshot of host resources.
//   - **`ctx system check_resource`** — the hook that fires a
//     pressure warning during sessions when load, memory, or
//     disk crosses a danger threshold.
//
// # Per-Platform Implementations
//
// Resource collection is **platform-conditional** via Go build
// tags so the binary stays a single static cross-compile while
// still asking each OS in its native dialect:
//
//   - **Linux** — reads `/proc/meminfo` and `/proc/loadavg`
//     directly ([memory_linux.go], [load_linux.go]).
//   - **macOS / Darwin** — shells out to `sysctl -n vm.loadavg`
//     and `vm_stat` and parses their output
//     ([memory_darwin.go], [load_darwin.go]).
//   - **Other / Windows** — stubs that return
//     `Supported: false` ([memory_other.go], [load_other.go],
//     [disk_windows.go]). The hook degrades gracefully rather
//     than aborting the session.
//
// Disk usage is read uniformly via `syscall.Statfs` on
// Unix-likes ([disk.go]) and stubbed on Windows.
//
// # Threshold Evaluation
//
// [threshold.go] holds the WARNING / DANGER cutoffs and the
// per-metric evaluator that turns a raw measurement into a
// severity. Defaults reflect "headroom you almost certainly
// want": load averages compared against CPU count, memory
// available below a percentage, disk free below a percentage.
// The 5-minute load average — not the 1-minute — is used to
// avoid false positives from transient spikes (a deliberate
// behavior, see commit `5958e558`).
//
// # The Output Shape
//
// [Resource] ([types.go]) is the unified record emitted by
// each collector: kind, value, unit, threshold, severity,
// support flag. [calc.go] holds the per-metric arithmetic
// (percent free, ratio computations) with explicit
// zero-division guards.
//
// # Concurrency
//
// Each call to a collector is a one-shot syscall + parse;
// nothing is cached at the package level. Concurrent callers
// produce independent readings. The `vm_stat` / `sysctl`
// shell-out path on macOS uses an external process which is
// the slowest case (~tens of milliseconds).
//
// # Related Packages
//
//   - [internal/cli/sysinfo]                    — the
//     top-level `ctx sysinfo` CLI surface.
//   - [internal/cli/system/cmd/check_resource]  — the
//     resource-pressure hook.
//   - [internal/config/sysinfo]                 — procfs
//     paths and meminfo field-key constants used by Linux.
//   - [internal/write/sysinfo]                  — terminal
//     output formatters consumed by both the CLI and hook.
package sysinfo
