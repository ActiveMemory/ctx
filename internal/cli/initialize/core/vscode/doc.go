//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package vscode generates VS Code workspace configuration
// files during the ctx init pipeline.
//
// # Overview
//
// When ctx init runs, this package creates the .vscode/
// directory and populates it with configuration files
// that integrate VS Code with ctx. Each file generator
// is idempotent: existing files are skipped to preserve
// user customisations.
//
// # Behavior
//
// createVSCodeArtifacts generates extensions.json (extension
// recommendations), tasks.json (ctx shell tasks), and mcp.json
// (MCP server registration) inside .vscode/, skipping files
// that already exist.
//
// # Generated Files
//
// The package creates three files in .vscode/:
//
//   - extensions.json -- recommends the ctx VS Code
//     extension. If the file exists and already lists
//     the extension, it is skipped. If the file exists
//     without the recommendation, the user is prompted
//     to add it manually.
//   - tasks.json -- defines shell tasks for common ctx
//     commands (status, drift, agent). Uses VS Code
//     task schema version 2.0.0 with shared terminal
//     panels.
//   - mcp.json -- registers the ctx MCP server so
//     VS Code can communicate with ctx via the Model
//     Context Protocol.
//
// # Internal Types
//
//   - [vsTask] -- single task definition for tasks.json
//   - [vsPresentation] -- terminal display settings
//   - [vsTasksFile] -- top-level tasks.json structure
//   - [vsMCPServer] -- MCP server entry in mcp.json
//   - [vsMCPFile] -- top-level mcp.json structure
//
// Individual file errors are non-fatal and reported
// inline, allowing the rest of the init pipeline to
// continue.
package vscode
