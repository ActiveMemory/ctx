//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// padEntryIDPattern matches a stable ID prefix like "[17] ".
const padEntryIDPattern = `^\[(\d+)\]\s`

// PadEntryID matches a stable ID prefix on scratchpad entries.
//
// Groups:
//   - 1: numeric ID
var PadEntryID = regexp.MustCompile(padEntryIDPattern)
