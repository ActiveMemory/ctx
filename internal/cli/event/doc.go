//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package event implements the ctx hook event command.
//
// Queries the local hook event log and displays matching
// entries. Events are recorded by system hooks during AI
// sessions and capture hook firings, session lifecycle
// transitions, and notification deliveries.
//
// # Filtering
//
// Results can be narrowed by hook name (--hook), session
// ID (--session), and event type (--event). The --last
// flag limits the number of returned entries (default
// from [config/event.DefaultLast]). The --all flag
// includes rotated log files that are normally excluded.
//
// # Output Formats
//
// Human-readable output shows a table with timestamp,
// hook, event, and session columns. JSON output (--json)
// emits an array of event objects for scripting.
//
// [Cmd] builds the cobra command with filter and format flags.
// [Run] queries the event log, applies filters, and renders
// matching entries as a table or JSON array.
package event
