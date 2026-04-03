//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"encoding/json"
	"net/url"
	"path/filepath"
	"runtime"

	cfgCopilot "github.com/ActiveMemory/ctx/internal/config/copilot"
	"github.com/ActiveMemory/ctx/internal/config/env"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// resolveWorkspaceCWD reads workspace.json from the workspaceStorage
// directory to determine the workspace folder path.
//
// Parameters:
//   - sessionPath: path to the JSONL session file
//
// Returns:
//   - string: the resolved workspace folder path, or empty string on failure
func (p *Copilot) resolveWorkspaceCWD(sessionPath string) string {
	// sessionPath is like: .../workspaceStorage/<hash>/chatSessions/<id>.jsonl
	// workspace.json is at: .../workspaceStorage/<hash>/workspace.json
	chatDir := filepath.Dir(sessionPath) // chatSessions/
	storageDir := filepath.Dir(chatDir)  // <hash>/
	wsFile := filepath.Join(storageDir, cfgCopilot.FileWorkspace)

	data, readErr := ctxIo.SafeReadUserFile(
		filepath.Clean(wsFile),
	)
	if readErr != nil {
		return ""
	}

	var ws copilotRawWorkspace
	if unmarshalErr := json.Unmarshal(
		data, &ws,
	); unmarshalErr != nil {
		return ""
	}

	return fileURIToPath(ws.Folder)
}

// fileURIToPath converts a file:// URI to a local file path.
//
// Parameters:
//   - uri: the file URI to convert (e.g., "file:///home/user/project")
//
// Returns:
//   - string: the local file path, or empty string if the URI is invalid
func fileURIToPath(uri string) string {
	if uri == "" {
		return ""
	}

	parsed, parseErr := url.Parse(uri)
	if parseErr != nil {
		return ""
	}

	if parsed.Scheme != cfgCopilot.SchemeFile {
		return ""
	}

	path := parsed.Path

	// URL-decode the path (e.g., %3A -> :)
	decoded, unescapeErr := url.PathUnescape(path)
	if unescapeErr != nil {
		decoded = path
	}

	// On Windows, file URIs have /G:/... — strip the leading slash
	if runtime.GOOS == env.OSWindows && len(decoded) > 2 && decoded[0] == '/' {
		decoded = decoded[1:]
	}

	return filepath.FromSlash(decoded)
}
