//go:build windows

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import "os"

func collectDisk() DiskInfo {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return DiskInfo{Supported: false}
	}
	return DiskInfo{Path: cwd, Supported: false}
}
