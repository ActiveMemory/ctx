//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package vscode generates VS Code workspace configuration files
// during ctx init.
//
// [WriteAll] is the entry point, invoked from the init pipeline.
// It delegates to per-file generators that create:
//
//   - extensions.json — recommended extensions including the ctx
//     VS Code extension
//   - tasks.json — shell tasks for common ctx commands (status,
//     drift, agent)
//   - mcp.json — MCP server registration pointing at ctx mcp serve
//
// Each generator skips its file if it already exists, printing a
// diagnostic via [writeVscode.InfoExistsSkipped]. Types used for
// JSON serialisation live in types.go.
package vscode
