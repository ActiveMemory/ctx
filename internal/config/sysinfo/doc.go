//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sysinfo defines constants for cross-platform
// system information collection used by ctx system
// bootstrap and the doctor command.
//
// ctx reports host resource data (memory, swap, load,
// disk) so agents understand the machine they run on.
// This package centralizes the platform-specific
// parsing vocabulary for both Linux and macOS.
//
// # Linux procfs Constants
//
//   - [ProcLoadavg], [ProcMeminfo]: file paths in
//     /proc/ for load averages and memory stats.
//   - [LoadavgFmt]: scanf format for three float
//     load averages.
//   - [FieldMemTotal], [FieldMemAvailable],
//     [FieldMemFree], [FieldBuffers], [FieldCached],
//     [FieldSwapTotal], [FieldSwapFree]: keys for
//     parsing /proc/meminfo lines.
//   - [BytesPerKB]: unit conversion factor.
//
// # macOS Constants
//
//   - [CmdSysctl], [CmdVMStat]: system commands.
//   - [KeyLoadAvg], [KeyHWMemsize],
//     [KeyVMSwapUsage]: sysctl keys for load,
//     memory, and swap.
//   - [MarkerPageSize], [LabelPagesFree],
//     [LabelPagesInactive]: vm_stat output parsing.
//   - [SuffixMB], [LabelTotal], [LabelUsed]: swap
//     usage parsing tokens.
//
// # Severity Labels
//
//   - [LabelOK], [LabelWarning], [LabelDanger]:
//     severity strings for resource threshold
//     evaluation.
//
// # Resource Names
//
//   - [ResourceMemory], [ResourceSwap],
//     [ResourceDisk], [ResourceLoad]: identifiers
//     for threshold lookup.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package sysinfo
