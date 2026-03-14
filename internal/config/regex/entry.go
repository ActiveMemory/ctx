//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// EntryHeader matches entry headers like "## [2026-01-28-051426] Title here".
//
// Groups:
//   - 1: date (YYYY-MM-DD)
//   - 2: time (HHMMSS)
//   - 3: title
var EntryHeader = regexp.MustCompile(
	`## \[(\d{4}-\d{2}-\d{2})-(\d{6})] (.+)`,
)

// EntryHeaderGroups is the expected number of groups (including full
// match) returned by EntryHeader.FindStringSubmatch.
const EntryHeaderGroups = 4

// EntryHeading matches any entry heading (## [timestamp]).
// Use for counting entries without capturing groups.
var EntryHeading = regexp.MustCompile(`(?m)^## \[`)
