//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package fs

// File permission constants.
const (
	// PermFile is the standard permission for regular files (owner rw, others r).
	PermFile = 0644
	// PermExec is the standard permission for directories and executable files.
	PermExec = 0755
	// PermRestrictedDir is the permission for internal
	// directories (owner rwx, group rx).
	PermRestrictedDir = 0750
	// PermSecret is the permission for secret files (owner rw only).
	PermSecret = 0600
	// PermKeyDir is the permission for the user-level key
	// directory (owner rwx only).
	PermKeyDir = 0700
)
