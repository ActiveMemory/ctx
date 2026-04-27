//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	cfgDir "github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	mcpServer "github.com/ActiveMemory/ctx/internal/config/mcp/server"
	cfgShell "github.com/ActiveMemory/ctx/internal/config/shell"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// launchCommand returns the OpenCode `command` array for the ctx
// MCP server. The emitted argv is:
//
//	["sh", "-c", `exec env CTX_DIR="$PWD/.context" ctx mcp serve`]
//
// `$PWD` is set by sh to the CWD OpenCode chose when spawning the
// MCP child — the project root that owns opencode.json. `exec`
// replaces the shell so the MCP child becomes ctx itself, with no
// sh process layered between OpenCode and the JSON-RPC stream.
//
// Returns:
//   - []string: argv suitable for OpenCode's McpLocalConfig.command.
func launchCommand() []string {
	binAndArgs := append([]string{mcpServer.Command}, mcpServer.Args()...)
	script := fmt.Sprintf(
		cfgShell.FormatPOSIXSpawnRelativeCtxDir,
		env.CtxDir, cfgDir.Context,
		strings.Join(binAndArgs, token.Space),
	)
	return []string{cfgShell.Sh, cfgShell.CmdFlag, script}
}

// ensureMCPConfig registers the ctx MCP server in opencode.json
// at the project root.
//
// Merge-safe: reads existing config, adds ctx server under
// the "mcp" key, writes back. Skips if ctx server is already
// registered. Treats a missing or empty file as "no existing
// config" rather than an error.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if file read/write fails
func ensureMCPConfig(cmd *cobra.Command) error {
	target := cfgHook.FileOpenCodeJSON

	// Read existing config if it exists. An empty or whitespace-only
	// file is treated as "no existing config" so users who pre-create
	// opencode.json don't trip an unmarshal error.
	existing := make(map[string]interface{})
	data, readErr := ctxIo.SafeReadUserFile(target)
	if readErr == nil && len(bytes.TrimSpace(data)) > 0 {
		if jErr := json.Unmarshal(data, &existing); jErr != nil {
			return jErr
		}
	}

	// Get or create mcp map.
	servers, _ := existing[cfgHook.KeyMCP].(map[string]interface{})
	if servers == nil {
		servers = make(map[string]interface{})
	}

	// Check if ctx is already registered.
	if _, ok := servers[mcpServer.Name]; ok {
		writeSetup.InfoOpenCodeSkipped(cmd, target)
		return nil
	}

	// Add ctx MCP server. OpenCode's McpLocalConfig schema differs
	// from Copilot CLI's: `command` is an Array<string> that holds
	// both the binary and its args (no separate `args` field), and
	// `enabled` is required at runtime.
	//
	// We wrap the launch in `sh -c` and resolve CTX_DIR to
	// `$PWD/.context` at spawn time. ctx requires CTX_DIR to be
	// absolute (see internal/rc.ContextDir), so a literal ".context"
	// in the `environment` map is rejected before the JSON-RPC
	// handshake. OpenCode also has no path templating in
	// opencode.json, so we cannot embed an absolute path that
	// follows the user's checkout. Computing it from $PWD at launch
	// gives us an absolute path anchored to the project that owns
	// this opencode.json, regardless of the user's shell CTX_DIR.
	servers[mcpServer.Name] = map[string]interface{}{
		cfgHook.KeyType:    cfgHook.MCPServerType,
		cfgHook.KeyCommand: launchCommand(),
		cfgHook.KeyEnabled: true,
	}
	existing[cfgHook.KeyMCP] = servers

	out, marshalErr := json.MarshalIndent(
		existing, "", token.Indent2,
	)
	if marshalErr != nil {
		return marshalErr
	}
	out = append(out, token.NewlineLF...)

	writeFileErr := ctxIo.SafeWriteFile(
		target, out, fs.PermFile,
	)
	if writeFileErr != nil {
		return writeFileErr
	}
	writeSetup.InfoOpenCodeCreated(cmd, target)
	return nil
}
