//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package gio wraps GNOME GIO command execution.
//
// Used for mounting SMB shares via gio mount during backup
// operations. The mount target URL comes from user configuration.
//
// Key exports: [Mount].
// Part of the exec subsystem.
package gio
