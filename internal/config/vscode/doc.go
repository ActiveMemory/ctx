//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package vscode defines constants for generating
// VS Code workspace configuration files during
// ctx init.
//
// When ctx initializes a project, it can scaffold a
// .vscode/ directory with extensions.json (to
// recommend the ctx extension), tasks.json (to
// expose common ctx commands as VS Code tasks), and
// mcp.json (to configure MCP server integration).
// This package provides every directory path, file
// name, JSON key, and task definition needed.
//
// # Directory and File Paths
//
//   - [Dir] (".vscode") — the workspace config dir.
//   - [FileExtensionsJSON], [FileTasksJSON],
//     [FileMCPJSON] — config file names within it.
//
// # Extension
//
//   - [ExtensionID] ("activememory.ctx-context") —
//     the VS Code Marketplace identifier for the
//     ctx extension.
//
// # JSON Keys
//
//   - [KeyRecommendations] — extensions.json key.
//   - [KeyCommand] — tasks.json command key.
//   - [KeyServers], [KeyArgs] — mcp.json keys.
//
// # Task Configuration
//
//   - [TasksVersion] ("2.0.0") — tasks schema
//     version.
//   - [TypeShell], [GroupNone], [RevealAlways],
//     [PanelShared] — task runner settings.
//   - [Tasks] — the label/command pairs written
//     into tasks.json (status, drift, agent,
//     journal, journal-serve).
//
// # Concurrency
//
// All exports are immutable except [Tasks], which
// is a package-level var but never mutated after
// init.
package vscode
