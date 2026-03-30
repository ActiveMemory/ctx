//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package index provides session ID indexing for journal files.
//
// [SessionIndex] builds a map from session IDs to filenames by
// scanning journal markdown frontmatter. [ExtractSessionID] pulls
// the session_id from a single file. [LookupSessionFile] resolves
// a session ID to its filename.
package index
