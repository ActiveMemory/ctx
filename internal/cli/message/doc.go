//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package message provides the ctx hook message command
// for injecting messages into AI sessions via the Claude
// Code hook message protocol.
//
// The message command reads structured input and produces
// hook-compatible JSON output that Claude Code interprets
// as injected context. This enables skills, triggers, and
// automation scripts to surface information to the AI
// agent mid-session without user intervention.
//
// # How It Works
//
// When a Claude Code hook fires, message reads the hook
// payload from stdin, evaluates whether a message should
// be injected, and writes a JSON response to stdout. The
// response may contain a message string that Claude Code
// prepends to the next AI prompt.
//
// # Subcommands
//
//   - edit: edit a previously sent message
//   - list: list messages in the current session
//   - reset: reset message state
//   - show: display a specific message
//
// # Subpackages
//
//	cmd/root: cobra command definition, stdin reading,
//	  and JSON response formatting
package message
