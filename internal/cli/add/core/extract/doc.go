//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package extract reads entry content from the available
// input sources for the add command.
//
// # Content Resolution
//
// [Content] resolves the entry body by checking three
// sources in strict priority order:
//
//  1. File flag: when AddConfig.FromFile is set, the
//     file is read via io.SafeReadUserFile and its
//     trimmed contents are returned.
//  2. Positional arguments: when args has more than
//     one element, args[1:] are joined with a space
//     separator.
//  3. Piped stdin: when stdin is not a terminal
//     (character device), all lines are read with a
//     bufio.Scanner and joined with newlines.
//
// If none of the sources produce content, Content returns
// an errAdd.NoContent error so the cmd/ layer can display
// a usage hint.
//
// # Error Handling
//
// File read failures surface as errFs.FileRead errors
// that include the original path. Stdin scanner errors
// surface as errFs.StdinRead. Both wrap the underlying
// OS error for inspection.
//
// # Data Flow
//
// The cmd/ layer calls Content early in the add pipeline.
// The returned string is then passed to a formatter in
// the format subpackage and finally to insert.AppendEntry
// for placement into the target context file.
package extract
