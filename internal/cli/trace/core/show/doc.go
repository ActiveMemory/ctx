//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package show implements the core logic for displaying commit context
// traces.
//
// [Commit] resolves refs for a single commit and prints them as text
// or JSON. [Last] iterates the last N commits from git log and
// summarises each with its context refs. [ResolveToJSON] converts raw
// refs to [JSONRef] structs for structured output.
//
// Key exports: [Commit], [Last], [ResolveToJSON], [JSONCommit], [JSONRef].
// Called by the trace show CLI subcommand.
package show
