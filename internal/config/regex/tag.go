//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// Tag matches #word tokens bounded by whitespace or start of string.
//
// Groups:
//   - 1: tag name (letters, digits, hyphens, underscores)
var Tag = regexp.MustCompile(`(?:^|\s)#([a-zA-Z0-9_-]+)`)
