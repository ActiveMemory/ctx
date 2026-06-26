//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package setup

import "github.com/ActiveMemory/ctx/internal/config/token"

// Root setup file paths.
const (
	// FileCtxRC is the project rc file name.
	FileCtxRC = ".ctxrc"
)

// Backend setup output strings.
const (
	BackendDryRunPrefix        = "Add this backend configuration to .ctxrc:\n"
	BackendEndpointRequired    = "%s endpoint is required"
	BackendEndpointScheme      = "%s endpoint must be http or https"
	BackendEnvAnthropicBaseURL = "ANTHROPIC_BASE_URL"
	BackendEnvOpenAIBaseURL    = "OPENAI_BASE_URL"
	BackendEnvWarn             = "warning: environment variable already set: "
	BackendEnvWarnEnd          = token.NewlineLF
	BackendRCMapping           = ".ctxrc must be a mapping"
	BackendShellExportLine     = "export %s=%q\n"
	BackendUnsupported         = "unsupported backend: "
	BackendWriteDone           = "Updated .ctxrc backend configuration\n"
)

// Display names for supported integration tools.
const (
	// DisplayKiro is the display name for Kiro.
	DisplayKiro = "Kiro"
	// DisplayCursor is the display name for Cursor.
	DisplayCursor = "Cursor"
	// DisplayCline is the display name for Cline.
	DisplayCline = "Cline"
)

// Kiro configuration paths.
const (
	// DirKiro is the Kiro editor config directory.
	DirKiro = ".kiro"
	// DirSettings is the Kiro settings subdirectory.
	DirSettings = "settings"
	// FileMCPJSON is the Kiro MCP server config file name.
	FileMCPJSON = "mcp.json"
	// MCPConfigPathKiro is the deployed MCP config path.
	MCPConfigPathKiro = DirKiro + "/settings/mcp.json"
	// SteeringDeployPathKiro is the deployed steering
	// path for Kiro.
	SteeringDeployPathKiro = DirKiro + "/steering/"
)

// Cursor configuration paths.
const (
	// DirCursor is the Cursor editor config directory.
	DirCursor = ".cursor"
	// FileMCPJSONCursor is the Cursor MCP config file.
	FileMCPJSONCursor = "mcp.json"
	// MCPConfigPathCursor is the deployed MCP config path.
	MCPConfigPathCursor = DirCursor + "/mcp.json"
	// SteeringPathCursor is the deployed steering path
	// for Cursor.
	SteeringPathCursor = DirCursor + "/rules/"
)

// OpenCode configuration paths.
const (
	// MCPConfigPathOpenCode is the OpenCode MCP config path used as a
	// fallback in user-facing warnings when globalConfigPath() cannot
	// resolve the actual path (e.g. UserHomeDir failure).
	MCPConfigPathOpenCode = "~/.config/opencode/opencode.json"
)

// Cline configuration paths.
const (
	// MCPConfigPathCline is the deployed MCP config path.
	MCPConfigPathCline = ".vscode/mcp.json"
	// SteeringPathCline is the deployed steering path
	// for Cline.
	SteeringPathCline = ".clinerules/"
)
