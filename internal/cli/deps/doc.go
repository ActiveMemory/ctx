//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package deps provides the `ctx deps` command for generating
// dependency graphs from source code. Supports Go, Node.js, Python,
// and Rust projects via ecosystem-specific GraphBuilder implementations.
// The ecosystem is auto-detected from manifest files or forced with --type.
package deps
