//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reduce strips formatting artifacts from session JSONL
// for clean journal markdown.
//
// [StripFences] removes code fence markers. [StripSystemReminders]
// removes system-reminder XML tags. [CleanToolOutputJSON] simplifies
// tool output JSON for readability.
package reduce
