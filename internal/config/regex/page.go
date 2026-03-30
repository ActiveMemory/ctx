//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// MultiPart matches session part files like "...-p2.md", "...-p3.md", etc.
var MultiPart = regexp.MustCompile(`-p\d+\.md$`)

// GlobStar matches glob-like wildcards: *.ext, */, *) etc.
var GlobStar = regexp.MustCompile(`\*(\.\w+|[/)])`)

// BlogDateFilename matches blog post filenames: YYYY-MM-DD-*.md.
var BlogDateFilename = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-.+\.md$`)
