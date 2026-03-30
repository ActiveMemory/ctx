//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core contains dependency graph building and rendering
// for Go, Node.js, and Rust ecosystems.
//
// [DetectBuilder] auto-selects the right builder based on project
// files. [FindBuilder] looks up by name. Each ecosystem has a
// GraphBuilder implementation that produces a directed graph.
// [MermaidID] sanitizes package names for Mermaid diagram syntax.
package core
