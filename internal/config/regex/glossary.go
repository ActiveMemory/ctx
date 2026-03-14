//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// Glossary matches glossary definition entries (lines with **term**).
var Glossary = regexp.MustCompile(`(?m)(?:^|\n)\s*(?:-\s*)?\*\*[^*]+\*\*`)
