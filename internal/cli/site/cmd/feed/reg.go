//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package feed

import "regexp"

// regBlogDatePattern matches filenames like 2026-02-25-slug.md.
var regBlogDatePattern = regexp.MustCompile(
	`^\d{4}-\d{2}-\d{2}-.+\.md$`,
)
