//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package marker

import "github.com/ActiveMemory/ctx/internal/config/token"

// Index block format strings for inserting/appending index content.
var (
	// IndexBlockFmt formats an index block between existing content.
	// Args: content-before, index-content, content-after.
	IndexBlockFmt = "%s" + token.NewlineLF +
		IndexStart + token.NewlineLF +
		"%s" + IndexEnd + token.NewlineLF + "%s"

	// IndexBlockAppendFmt appends an index block at end of file.
	// Args: content, index-content.
	IndexBlockAppendFmt = "%s" + token.NewlineLF + token.NewlineLF +
		IndexStart + token.NewlineLF +
		"%s" + IndexEnd + token.NewlineLF
)
