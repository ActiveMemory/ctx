//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cooldown prevents redundant context loads within a
// single session.
//
// [Active] checks whether a cooldown tombstone exists and is
// recent enough. [TouchTombstone] creates or refreshes the
// tombstone file. [TombstonePath] returns the state file path
// for a given session ID.
package cooldown
