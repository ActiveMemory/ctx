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
//   - [ProcPressureMemory], [PSILineSome],
//     [PSILineFull], [PSIFieldAvg10]: path and tokens
//     for parsing the /proc/pressure/memory PSI signal.
//   - [BytesPerKB]: unit conversion factor.
//
// # macOS Constants
//
//   - [CmdSysctl], [CmdVMStat]: system commands.
//   - [KeyLoadAvg], [KeyHWMemsize],
//     [KeyVMSwapUsage], [KeyVMPressureLevel]: sysctl
//     keys for load, memory, swap, and the memory
//     pressure level.
//   - [PressureLevelNormal], [PressureLevelWarning],
//     [PressureLevelCritical]: kern.memorystatus_vm_pressure_level
//     values.
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
//   - [ResourceMemory], [ResourceMemoryPressure],
//     [ResourceSwap], [ResourceDisk], [ResourceLoad]:
//     identifiers for threshold lookup.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package sysinfo
