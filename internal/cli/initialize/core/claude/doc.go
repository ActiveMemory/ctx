//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package claude handles CLAUDE.md creation and merging during init.
//
// [HandleMd] creates CLAUDE.md if missing, or merges the ctx
// section into an existing file using marker-delimited regions.
// Force mode overwrites without merging.
package claude
