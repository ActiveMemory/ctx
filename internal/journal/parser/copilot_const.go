//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

// Copilot JSONL line Kind values.
const (
	// copilotKindSnapshot is a full session snapshot (kind=0).
	copilotKindSnapshot = 0
	// copilotKindScalarPatch is a scalar field replacement (kind=1).
	copilotKindScalarPatch = 1
	// copilotKindObjectPatch is an array/object replacement (kind=2).
	copilotKindObjectPatch = 2
)

// Copilot JSON key paths and response item kinds.
const (
	copilotKeyRequests = "requests"
	copilotKeyResult   = "result"
	copilotKeyResponse = "response"

	copilotRespKindThinking   = "thinking"
	copilotRespKindToolInvoke = "toolInvocationSerialized"

	copilotDirChatSessions = "chatSessions"
	copilotFileWorkspace   = "workspace.json"
	copilotResponseSuffix  = "-response"
)

// Copilot scanner buffer sizes for JSONL parsing.
const (
	// copilotScanBufInit is the initial scanner buffer size (64KB).
	copilotScanBufInit = 64 * 1024
	// copilotScanBufMax is the maximum scanner buffer size (4MB).
	// Copilot lines can be very large due to embedded code content.
	copilotScanBufMax = 4 * 1024 * 1024
	// copilotScanBufMatchMax is the maximum scanner buffer for Matches
	// checks (1MB). Smaller than full parse because only the first line
	// is inspected.
	copilotScanBufMatchMax = 1024 * 1024
)

// Copilot tool ID and display constants.
const (
	// copilotToolIDSeparator separates the namespace prefix from the
	// tool name in Copilot tool IDs (e.g., "copilot_readFile").
	copilotToolIDSeparator = "_"
	// copilotCLIAppName is the application directory name used on
	// Windows under LOCALAPPDATA for Copilot CLI sessions.
	copilotCLIAppName = "GitHub Copilot CLI"
)

// Copilot platform and path constants used in CopilotSessionDirs.
const (
	copilotEnvAppData      = "APPDATA"
	copilotOSDarwin        = "darwin"
	copilotSchemeFile      = "file"
	copilotAppCode         = "Code"
	copilotAppCodeInsiders = "Code - Insiders"
	copilotDirUser         = "User"
	copilotDirWorkspace    = "workspaceStorage"
	copilotDirLibrary      = "Library"
	copilotDirAppSupport   = "Application Support"
	copilotDirDotConfig    = ".config"
)
