//go:build !windows

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import "syscall"

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
