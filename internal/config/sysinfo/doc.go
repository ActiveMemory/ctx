//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sysinfo provides constants for system information collection
// including Linux procfs paths and meminfo field keys.
//
// Procfs constants ([ProcLoadavg], [ProcMeminfo]) define file paths for
// reading load averages and memory statistics. Meminfo field keys
// ([FieldMemTotal], [FieldMemAvailable], etc.) identify lines in
// /proc/meminfo for parsing. [BytesPerKB] provides the unit conversion
// factor.
//
// Used by the system bootstrap command to report host resource data.
package sysinfo
