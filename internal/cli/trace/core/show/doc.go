//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package show displays context references attached to
// git commits. It supports both human-readable text and
// structured JSON output.
//
// # Single Commit Display
//
// [Commit] resolves context refs for a single commit
// hash. In text mode it prints a header with the short
// hash, message, and date, followed by each resolved
// ref with its type, number, and title. In JSON mode
// it encodes a [JSONCommit] struct containing the commit
// metadata and a list of [JSONRef] entries.
//
// # Bulk Listing
//
// [Last] iterates the last N commits from git log and
// summarizes each with its context refs. For efficiency,
// bulk listing skips trailer re-reading because the
// history file already contains the same refs the
// post-commit hook extracted from the trailer.
//
// # JSON Types
//
// [JSONCommit] pairs a commit hash and message with its
// resolved refs. [JSONRef] holds a single resolved
// reference with its raw string, type, optional number
// and title, detail text, and a found flag indicating
// whether the ref could be resolved.
//
// # Reference Resolution
//
// [ResolveToJSON] converts raw ref strings to JSONRef
// structs by calling trace.Resolve against the context
// directory.
package show
