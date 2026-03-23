//go:build !windows

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import (
	"os"
	"syscall"
)

// collectDisk queries filesystem statistics for the current working directory.
//
// Uses syscall.Statfs to obtain total and available block counts,
// then converts to byte values. Returns a DiskInfo with Supported=false
// if the working directory cannot be determined or statfs fails.
//
// Returns:
//   - DiskInfo: Disk usage statistics for the filesystem containing CWD
func collectDisk() DiskInfo {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return DiskInfo{Supported: false, Err: cwdErr}
	}

	var stat syscall.Statfs_t

	if statErr := syscall.Statfs(cwd, &stat); statErr != nil {
		return DiskInfo{Path: cwd, Supported: false, Err: statErr}
	}
	if stat.Bsize <= 0 {
		return DiskInfo{Path: cwd, Supported: false}
	}
	bsize := uint64(stat.Bsize)
	total := stat.Blocks * bsize
	free := stat.Bavail * bsize // available to unprivileged users
	var used uint64
	if total > free {
		used = total - free
	}
	return DiskInfo{
		TotalBytes: total,
		UsedBytes:  used,
		Path:       cwd,
		Supported:  true,
	}
}
