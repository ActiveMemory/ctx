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
//   - **[EntriesWithIDs](raw)** and
//     **[FormatEntriesWithIDs](entries)**: the
//     ID-aware pair; parse lines carrying stable
//     `[N] ` prefixes into `[]Entry` and serialize
//     them back to the on-disk shape.
//
// # Round-Trip Stability
//
// `FormatEntriesWithIDs(EntriesWithIDs(x))` is
// byte-identical to `x` when every line of `x`
// already carries a unique `[N] ` prefix. Lines
// without a prefix (or with a duplicate ID) are
// assigned fresh IDs on parse, so the first write
// normalizes them.
//
// # Concurrency
//
// Pure data transformation. Concurrent callers
// never race.
package parse
