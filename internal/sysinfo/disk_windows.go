//go:build windows

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import "os"

// collectDisk returns disk usage information for the current working directory.
// On Windows this is a stub that always reports Supported as false.
//
// Returns:
//   - DiskInfo: Disk info with Supported set to false
func collectDisk() DiskInfo {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return DiskInfo{Supported: false}
	}
	return DiskInfo{Path: cwd, Supported: false}
}
