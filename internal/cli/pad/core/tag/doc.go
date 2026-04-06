//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package tag extracts and matches #word tags in scratchpad entries.
//
// Key exports: [Extract], [Has], [Match], [MatchAll], [ScanText].
// Tags are convention-based: any #word token in entry text is a tag.
// Used by the pad root command for filtering and the tags subcommand.
package tag
