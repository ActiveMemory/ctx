//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx watch" cobra
// command.
//
// This command reads a stream of context update
// commands from a log file or stdin and applies them
// to the project's context files. It is designed for
// processing structured output from agent sessions
// in real time.
//
// # Usage
//
//	ctx watch [--log <file>] [--dry-run]
//
// # Flags
//
//	--log       Path to a log file to read. When
//	            omitted, reads from stdin.
//	--dry-run   Show what updates would be applied
//	            without actually modifying any
//	            context files.
//
// # Behavior
//
// The command:
//
//   - Verifies that the context directory exists;
//     returns an error if it has not been initialized.
//   - Prints a startup banner with a hint on how to
//     stop the watcher (Ctrl-C).
//   - Opens the log file (or uses stdin) as the
//     input reader.
//   - Delegates to watch/core/stream.Process which
//     reads lines from the stream, parses context
//     update commands, and applies them to the
//     appropriate context files.
//   - In dry-run mode, stream.Process prints what
//     it would do without writing to disk.
//
// # Output
//
// Prints a startup banner, then live status lines as
// context updates are processed. In dry-run mode,
// each update is prefixed with a dry-run indicator.
//
// # Delegation
//
// Stream processing is in watch/core/stream. Context
// validation uses context/validate. Output formatting
// uses write/watch.
package root
