//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package imp

// BlobResult holds the outcome of importing a single file as a blob.
type BlobResult struct {
	Name     string
	Err      error
	TooLarge bool
	Added    bool
}
