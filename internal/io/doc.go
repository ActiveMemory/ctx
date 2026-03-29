//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package line provides guarded file I/O wrappers for ctx.
//
// # What these functions guard against
//
// All Safe* functions apply two checks before touching the filesystem:
//
//   - Path cleaning: filepath.Clean removes redundant separators,
//     dot segments, and trailing slashes.
//   - System prefix rejection: the resolved absolute path is checked
//     against a deny list of system directories (/bin, /etc, /proc,
//     /sys, /dev, /boot, /lib, /sbin, /usr/bin, /usr/lib, /usr/sbin,
//     and the filesystem root itself). Any match returns an error
//     before the underlying syscall executes.
//
// SafeReadFile additionally enforces containment: the resolved path
// must stay within the provided base directory.
//
// # What these functions do NOT guard against
//
//   - Symlink attacks: a cleaned path that passes the prefix check
//     could still resolve to a different location through a symlink
//     in a parent directory. Use [validation.CheckSymlinks] separately
//     when the directory tree is untrusted.
//   - Race conditions (TOCTOU): the check and the I/O are not atomic.
//     A malicious actor with write access to the parent directory
//     could swap a path between validation and use.
//   - Permission escalation: these wrappers run with the calling
//     process's permissions. They do not drop privileges.
//   - Content validation: the wrappers check where data is read from
//     or written to, not what the data contains.
//   - Windows paths: the deny list uses Unix prefixes. On Windows,
//     the prefix check is effectively a no-op.
//
// # Assumptions
//
// Callers are expected to provide paths that are already logically
// correct (e.g., constructed from known config constants or user
// input that has been validated for format). These wrappers are a
// safety net against accidental system directory access, not a
// substitute for input validation at the application boundary.
//
// # When to use which function
//
//   - SafeReadFile: path is base + filename (boundary-checked read)
//   - SafeReadUserFile: single path from any source (deny-list read)
//   - SafeOpenUserFile: single path, need a file handle (deny-list open)
//   - SafeWriteFile: single path (deny-list write)
package io
