//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package copilot centralizes constants for GitHub Copilot
// Chat and Copilot CLI session parsing and integration.
//
// ctx imports session history from Copilot Chat (VS Code
// extension) and Copilot CLI into its journal. This
// package defines the JSON keys, response kinds, scanner
// limits, tool ID parsing rules, and platform-specific
// storage paths that the importers need.
//
// # JSON Key Paths
//
// Copilot Chat sessions are stored as JSONL. The importer
// navigates nested JSON using these keys:
//
//   - [KeyRequests] -- top-level array of chat turns
//   - [KeyResult] -- result object within a request
//   - [KeyResponse] -- response payload within a result
//
// # Response Kinds
//
// Each response item has a kind field. The importer
// filters on:
//
//   - [RespKindThinking] -- extended thinking blocks
//   - [RespKindToolInvoke] -- serialized tool invocations
//
// # Scanner Buffer Sizes
//
// Copilot JSONL lines can be very large (embedded code).
// The scanner buffers are sized accordingly:
//
//   - [ScanBufInit] -- 64KB initial buffer
//   - [ScanBufMax] -- 4MB ceiling for full parsing
//   - [ScanBufMatchMax] -- 1MB ceiling for format
//     detection (only the first line is inspected)
//
// # Tool ID Parsing
//
// [ToolIDSeparator] ("_") splits the namespace prefix
// from the tool name in Copilot tool IDs (e.g.
// "copilot_readFile" becomes namespace "copilot" and
// tool "readFile").
//
// # Storage Paths
//
// Session files live in platform-specific directories.
// Constants cover both VS Code and Copilot CLI paths:
//
//   - [DirChatSessions] / [FileWorkspace] -- VS Code
//     workspace storage
//   - [CLIAppName] / [DirSessions] / [DirHistory] --
//     Copilot CLI session directories
//   - [OSDarwin] / [DirLibrary] / [DirAppSupport] --
//     macOS path components
//   - [EnvAppData] -- Windows environment variable
//
// # Why Centralized
//
// The journal importer, the session matcher, and the
// setup command all resolve Copilot session paths. A
// single package prevents platform-path drift and keeps
// buffer sizes consistent across parsing stages.
package copilot
