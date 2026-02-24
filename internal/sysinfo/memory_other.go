//go:build !linux && !darwin

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

func collectMemory() MemInfo {
	return MemInfo{Supported: false}
}
