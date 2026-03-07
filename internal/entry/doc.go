//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package entry provides the domain API for adding entries to context files.
// It owns the EntryParams type, validation, formatting, and write logic.
// CLI commands and external consumers (mcp, watch, memory) import this
// package instead of reaching into cli/add subpackages.
package entry
