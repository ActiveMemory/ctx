//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package io provides guarded file I/O and HTTP
// wrappers for ctx.
//
// # Filesystem Guards
//
// All Safe* functions apply two checks before
// touching the filesystem:
//
//   - Path cleaning: filepath.Clean removes redundant
//     separators, dot segments, and trailing slashes.
//   - System prefix rejection: the resolved absolute
//     path is checked against a deny list of system
//     directories (/bin, /etc, /proc, /sys, /dev,
//     /boot, /lib, /sbin, /usr/bin, /usr/lib,
//     /usr/sbin, and root itself). Any match returns
//     an error before the syscall executes.
//
// # File Operations
//
//   - [SafeReadFile] reads a file with containment:
//     the resolved path must stay within a base dir.
//   - [SafeReadUserFile] reads after deny-list check.
//   - [SafeOpenUserFile] opens for reading after
//     deny-list check (caller must close).
//   - [SafeAppendFile] opens in append mode, creating
//     the file if missing.
//   - [SafeCreateFile] creates or truncates a file.
//   - [SafeWriteFile] writes data after deny-list
//     check.
//   - [SafeMkdirAll] creates a directory tree after
//     deny-list check.
//   - [SafeStat] returns file info after deny-list
//     check.
//   - [TouchFile] creates or updates an empty marker
//     file (best-effort, errors logged).
//   - [AppendBytes] appends data in append mode with
//     best-effort error logging (for JSONL logs).
//
// # Formatted Output
//
//   - [SafeFprintf] writes formatted output to a
//     writer, logging errors to the warning sink.
//
// # HTTP
//
//   - [SafePost] sends an HTTP POST with scheme
//     validation (http/https only), redirect cap
//     (max 3), and caller-specified timeout.
//
// # Limitations
//
// These wrappers do not guard against symlink attacks,
// TOCTOU race conditions, permission escalation,
// content validation, or Windows paths. See the
// function-level documentation for details.
package io
