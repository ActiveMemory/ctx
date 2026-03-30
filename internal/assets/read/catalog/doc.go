//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package catalog lists available context template files from the
// embedded assets.
//
// [List] returns the names of all .context/ template files
// (TASKS.md, DECISIONS.md, etc.) available for deployment by
// ctx init. The list is derived from the embedded filesystem
// at compile time.
package catalog
