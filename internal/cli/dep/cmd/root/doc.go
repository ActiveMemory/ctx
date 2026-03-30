//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the ctx dep command for generating
// dependency graphs.
//
// [Cmd] builds the cobra.Command with --format, --full, and
// --builder flags. [Run] auto-detects the project ecosystem
// (Go, Node, Rust), builds the dependency graph, and renders
// it in the requested format (table, mermaid, json).
package root
