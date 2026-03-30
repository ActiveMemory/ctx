//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package status provides terminal output for the context status
// command (ctx status).
//
// [Header] renders the context directory path with file count and
// total token estimate. [FileItem] renders one context file with
// its token count and age; verbose mode adds the full path.
// [Activity] renders recent session activity as a summary list.
//
// Types [FileInfo] and [ActivityInfo] carry pre-computed display
// data so the write functions contain no business logic.
package status
