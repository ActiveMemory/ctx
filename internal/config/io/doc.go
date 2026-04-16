//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package io centralizes path-safety constants that
// govern where ctx is allowed to read and write on
// the filesystem.
//
// The primary export is [DangerousPrefixes], a slice
// of absolute directory prefixes (/bin, /boot, /dev,
// /etc, /proc, /sys, and others) that ctx must never
// touch. Every file-write path is resolved with
// filepath.Abs and checked against this list before
// any I/O proceeds.
//
// # Why Centralize
//
// Scattering safety checks across command handlers
// invites omissions. By defining the deny-list here,
// every caller shares the same boundary and any
// additions to the list take effect project-wide.
//
// # Key Constants
//
//   - DangerousPrefixes: system directories that ctx
//     must never read from or write to. Includes
//     /bin, /boot, /dev, /etc, /lib, /lib64, /proc,
//     /sbin, /sys, /usr/bin, /usr/lib, /usr/sbin.
//
// # How Consumers Use It
//
// A typical consumer resolves a user-supplied path to
// an absolute path, then iterates DangerousPrefixes
// to reject any match before opening the file. This
// prevents accidental or malicious writes to system
// directories during hook execution or file export.
package io
