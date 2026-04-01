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
