//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// CaseInsensitivePrefix is the inline flag prefix prepended to
// user-supplied regex patterns for case-insensitive matching.
const CaseInsensitivePrefix = "(?i)"

// CompileDynamic aliases regexp.Compile at package scope so callers
// can compile user-supplied CLI regex patterns through internal/config/regex.
var CompileDynamic = regexp.Compile
