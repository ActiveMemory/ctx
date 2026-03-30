//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the ctx agent command for generating
// AI-ready context packets.
//
// [Cmd] builds the cobra.Command with --budget, --format, and
// --json flags. [Run] loads context, assembles a budget-aware
// packet via core/budget, and renders it as Markdown or JSON.
package root
