//go:build windows

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

func collectDisk(path string) DiskInfo {
	return DiskInfo{Path: path, Supported: false}
}
