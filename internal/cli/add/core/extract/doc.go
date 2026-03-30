//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package extract reads entry content from command arguments or
// from a file specified by --from-file.
//
// [Content] joins positional arguments into a single string, or
// reads the file path from flags. Returns an error if both sources
// are empty or if the file cannot be read.
package extract
