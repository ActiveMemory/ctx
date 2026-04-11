//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

//go:build windows

package daemon

import "syscall"

// createNewProcessGroup is the Windows CREATE_NEW_PROCESS_GROUP
// creation flag. Declared as a constant rather than pulled from
// golang.org/x/sys/windows to avoid adding a dependency for a
// single value.
const createNewProcessGroup = 0x00000200

// detachAttrs returns the SysProcAttr used to detach a child
// process from the current console on Windows.
//
// CREATE_NEW_PROCESS_GROUP disables CTRL+C propagation from the
// parent console, letting the daemon outlive an interactive
// shell. HideWindow suppresses the console window flash for
// headless launches.
//
// Returns:
//   - *syscall.SysProcAttr: detached process-group attributes
func detachAttrs() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		CreationFlags: createNewProcessGroup,
		HideWindow:    true,
	}
}
