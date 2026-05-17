//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package closeout

import "github.com/ActiveMemory/ctx/internal/entity"

// Frontmatter aliases [entity.CloseoutFrontmatter] so the
// write/closeout package's existing call sites continue to read
// closeout.Frontmatter while the source-of-truth lives in
// entity/ (per the cross-package-types convention).
type Frontmatter = entity.CloseoutFrontmatter

// File aliases [entity.CloseoutFile] so the write/closeout
// package's existing call sites continue to read closeout.File
// while the source-of-truth lives in entity/.
type File = entity.CloseoutFile
