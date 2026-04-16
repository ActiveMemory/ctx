//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package daemon provides process management for
// background hub server operation.
//
// # Starting a Daemon
//
// Start launches a detached background process with
// the given binary path and arguments. It returns the
// PID of the started process. The child process is
// fully detached from the parent session so it
// survives when the parent shell exits.
//
//	pid, err := daemon.Start("/usr/bin/ctx", args)
//
// # Platform Detachment
//
// The detachAttrs function returns platform-specific
// SysProcAttr values:
//
//   - Unix: sets Setsid to create a new session,
//     making the child a session leader that is
//     independent of the parent terminal.
//   - Windows: sets CREATE_NEW_PROCESS_GROUP to
//     disable CTRL+C propagation and HideWindow
//     to suppress the console window flash.
//
// # Process Isolation
//
// The started process has nil stdout and stderr,
// ensuring it does not hold open the parent's file
// descriptors. This prevents the parent from hanging
// on exit while waiting for the child to close its
// output streams.
package daemon
