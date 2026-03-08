//go:build !linux && !darwin

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

// collectLoad is a no-op stub for unsupported platforms.
//
// Returns:
//   - LoadInfo: Always returns Supported=false
func collectLoad() LoadInfo {
	return LoadInfo{Supported: false}
}
