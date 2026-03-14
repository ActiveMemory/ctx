//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// LineNumber matches Claude Code's line number prefixes like "     1→".
var LineNumber = regexp.MustCompile(`(?m)^\s*\d+→`)
