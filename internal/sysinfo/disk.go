//go:build !windows

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import "syscall"

// collectDisk queries filesystem statistics for the given mount path.
//
// Uses syscall.Statfs to obtain total and available block counts,
// then converts to byte values. Returns a DiskInfo with Supported=false
// if the statfs call fails (e.g. path does not exist).
//
// Parameters:
//   - path: Filesystem path to query (typically "/" or a mount point)
//
// Returns:
//   - DiskInfo: Disk usage statistics for the filesystem containing path
func collectDisk(path string) DiskInfo {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return DiskInfo{Path: path, Supported: false}
	}
	bsize := uint64(stat.Bsize) //nolint:unconvert,gosec // type varies by OS; Bsize is always positive
	total := stat.Blocks * bsize
	free := stat.Bavail * bsize // available to unprivileged users
	var used uint64
	if total > free {
		used = total - free
	}
	return DiskInfo{
		TotalBytes: total,
		UsedBytes:  used,
		Path:       path,
		Supported:  true,
	}
}
