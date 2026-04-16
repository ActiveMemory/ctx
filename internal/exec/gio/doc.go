//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package gio wraps GNOME GIO command execution for
// mounting network shares.
//
// # Mounting
//
// Mount runs the "gio mount" command with the given
// URL. This is used during backup operations to mount
// SMB shares before copying context files.
//
//	err := gio.Mount("smb://host/share")
//
// # Dependencies
//
// The gio binary must be installed on the system.
// Mount returns an error if gio is not found in PATH
// or if the mount operation fails. The binary name
// and subcommand are sourced from the archive config
// package.
//
// # Security
//
// The URL argument comes from user configuration
// (not from untrusted input). The exec.Command call
// is annotated with a gosec nolint directive to
// acknowledge this.
package gio
