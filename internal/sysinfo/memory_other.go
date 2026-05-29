//go:build !linux && !darwin

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

// collectMemory is a no-op stub for unsupported platforms.
//
// Memory pressure is unsupported here, so PressureSupported is
// false and no pressure alert is ever raised on these platforms.
//
// Returns:
//   - MemInfo: Always returns Supported=false and PressureSupported=false
func collectMemory() MemInfo {
	return MemInfo{Supported: false, PressureSupported: false}
}
