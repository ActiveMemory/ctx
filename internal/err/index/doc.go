//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package index defines the typed error constructors
// returned by the index-block guard
// ([internal/index.Validate]).
//
// # Domain
//
// Every context file that carries an auto-generated
// index (DECISIONS.md, LEARNINGS.md) brackets that
// index between INDEX:START / INDEX:END markers. The
// index regenerator replaces the span between those
// markers wholesale. Before it runs, the guard checks
// that doing so is safe and refuses otherwise. Two
// refusals can occur:
//
//   - **Entries in block**: real entry bodies live
//     between the markers; regenerating would delete
//     them. Constructor: [EntriesInBlock].
//   - **Malformed markers**: the marker pair is
//     missing, duplicated, or out of order, so a
//     regenerate would emit a second marker.
//     Constructor: [MalformedMarkers].
//
// # Wrapping Strategy
//
// Both constructors return plain formatted errors;
// there is no underlying cause to wrap. The file name
// is interpolated so the message names the offending
// file. All user-facing text resolves through
// [internal/assets/read/desc] at construction time.
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package index
