//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// Slug matches characters that are NOT a-z, 0-9, or hyphen.
// Use it to strip a free-text title down to kebab-case-safe
// runes after lowercasing and replacing spaces with hyphens.
// Hyphens are allowed in the result; trim with `strings.Trim`.
var Slug = regexp.MustCompile(`[^a-z0-9-]+`)

// TopicSlug normalises free-text topic names into kebab-case
// slug-safe runes. Unlike [Slug], it preserves the `/`
// character so vendor-namespaced topic slugs (e.g.
// `cursor/hooks`) survive normalisation.
var TopicSlug = regexp.MustCompile(`[^a-z0-9/-]+`)
