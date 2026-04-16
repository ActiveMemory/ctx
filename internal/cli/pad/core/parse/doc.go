//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package parse splits raw scratchpad content into
// individual entries, the helper every `ctx pad`
// subcommand uses to turn the on-disk blob into a
// `[]Entry` it can filter, render, or mutate.
//
// # Public Surface
//
//   - **[Entries](raw)**: returns the entry slice
//     parsed from the scratchpad text. Recognizes
//     the `## YYYY-MM-DD HH:MM:SS` entry header;
//     everything between two headers (or between a
//     header and EOF) is one entry's body.
//   - **[FormatEntries](entries)**: the inverse;
//     serializes a `[]Entry` back to the raw on-disk
//     shape so writes round-trip cleanly.
//
// # Round-Trip Stability
//
// `FormatEntries(Entries(x))` is byte-identical to
// `x` when `x` is well-formed. This invariant is
// what makes `ctx pad edit` safe: the user's edits
// only land where the user typed.
//
// # Concurrency
//
// Pure data transformation. Concurrent callers
// never race.
package parse
