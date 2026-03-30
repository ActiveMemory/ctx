//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package collapse condenses large tool output sections in journal
// markdown.
//
// [ToolOutputs] finds tool output blocks in the content and
// replaces them with collapsed summaries, preserving the first
// few lines as context.
package collapse
