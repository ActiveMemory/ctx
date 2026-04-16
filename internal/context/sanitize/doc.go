//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sanitize provides content-level checks for
// context files.
//
// # Emptiness Detection
//
// EffectivelyEmpty checks whether a file contains only
// structural elements (headings, separators, HTML
// comment markers, whitespace) and no meaningful body
// content. This is used to filter out context files
// that exist on disk but have not yet been populated
// with real entries.
//
//	if sanitize.EffectivelyEmpty(content) {
//	    // skip this file in assembly
//	}
//
// # Heuristics
//
// The function applies several checks:
//
//   - Files shorter than content.MinLen are treated
//     as empty immediately.
//   - Lines starting with "#" are counted as headings
//     and skipped.
//   - Short lines starting with "-" are counted as
//     separators and skipped.
//   - Lines matching HTML comment open/close markers
//     are skipped.
//   - Any remaining non-whitespace line counts as
//     meaningful content.
//
// If zero content lines remain after filtering, the
// file is considered effectively empty.
package sanitize
