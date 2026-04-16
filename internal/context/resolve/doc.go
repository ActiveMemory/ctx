//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resolve provides helpers for locating the
// context directory and building derived paths.
//
// # Journal Directory
//
// JournalDir returns the absolute path to the journal
// subdirectory within the configured context directory.
// The context directory is read from the rc package.
//
//	path := resolve.JournalDir()
//	// => "/project/.context/journal"
//
// # Directory Line
//
// DirLine returns a one-line identifier string for the
// current context directory, formatted as a label. It
// returns an empty string if no directory is configured.
//
//	line := resolve.DirLine()
//	// => "Context: /project/.context"
//
// # Directory Footer
//
// AppendDir appends a bracketed context directory
// footer to a message string. If no context directory
// is available, the message is returned unchanged.
// This is used to annotate output messages with the
// active context path.
//
//	msg := resolve.AppendDir("Operation complete")
//	// => "Operation complete [Context: /path]"
package resolve
