//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package copilot defines constants for Copilot Chat and Copilot CLI
// session parsing and integration.
//
// Provides JSON key paths, response item kinds, scanner buffer sizes,
// VS Code storage paths, and Copilot CLI directory names used by the
// journal parser and setup command.
package copilot

// JSON key paths in Copilot Chat JSONL session files.
const (
	KeyRequests = "requests"
	KeyResult   = "result"
	KeyResponse = "response"
)

// Response item kind values in Copilot Chat sessions.
const (
	RespKindThinking   = "thinking"
	RespKindToolInvoke = "toolInvocationSerialized"
)

// Copilot Chat session storage directory and file names.
const (
	DirChatSessions = "chatSessions"
	FileWorkspace   = "workspace.json"
	ResponseSuffix  = "-response"
)

// Scanner buffer sizes for JSONL parsing.
const (
	// ScanBufInit is the initial scanner buffer size (64KB).
	ScanBufInit = 64 * 1024
	// ScanBufMax is the maximum scanner buffer size (4MB).
	// Copilot lines can be very large due to embedded code content.
	ScanBufMax = 4 * 1024 * 1024
	// ScanBufMatchMax is the maximum scanner buffer for Matches
	// checks (1MB). Smaller than full parse because only the first
	// line is inspected.
	ScanBufMatchMax = 1024 * 1024
)

// Tool ID parsing.
const (
	// ToolIDSeparator separates the namespace prefix from the tool
	// name in Copilot tool IDs (e.g., "copilot_readFile").
	ToolIDSeparator = "_"
)

// Copilot CLI application and session directory names.
const (
	// CLIAppName is the application directory name used on Windows
	// under LOCALAPPDATA for Copilot CLI sessions.
	CLIAppName = "GitHub Copilot CLI"
	// DirSessions is a candidate session subdirectory.
	DirSessions = "sessions"
	// DirHistory is a candidate session subdirectory.
	DirHistory = "history"
)

// VS Code platform and storage path constants.
const (
	EnvAppData      = "APPDATA"
	OSDarwin        = "darwin"
	SchemeFile      = "file"
	AppCode         = "Code"
	AppCodeInsiders = "Code - Insiders"
	DirUser         = "User"
	DirWorkspace    = "workspaceStorage"
	DirLibrary      = "Library"
	DirAppSupport   = "Application Support"
	DirDotConfig    = ".config"
)
