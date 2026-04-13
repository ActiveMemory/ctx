//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// FileNameChar matches characters not allowed in file names.
var FileNameChar = regexp.MustCompile(`[^a-zA-Z0-9-]+`)

// UUID matches a UUID (v4) anywhere in a string.
var UUID = regexp.MustCompile(
	`[0-9a-f]{8}-[0-9a-f]{4}-` +
		`[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`,
)

// SlugUnsafe matches characters replaced during slug
// generation (slash, dot). Glob star (*) is handled
// separately as it maps to a different replacement.
var SlugUnsafe = regexp.MustCompile(`[/.]`)
