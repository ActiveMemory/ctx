//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

//go:build !windows

package daemon

import "syscall"

// detachAttrs returns the SysProcAttr used to detach a child
// process from the current session on Unix-like systems.
//
// Setsid creates a new session so the child survives when the
// parent shell exits.
//
// Returns:
//   - *syscall.SysProcAttr: session-leader attributes
func detachAttrs() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setsid: true}
}
