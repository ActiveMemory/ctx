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
	"os"
	"os/exec"
	"path/filepath"
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
//	["sh", "-c", `exec env CTX_DIR="$PWD/.context" /abs/path/to/ctx mcp serve`]
//
// The binary path is resolved to an absolute path at setup time via
// exec.LookPath, so that OpenCode can spawn the MCP child regardless
// of the PATH inherited by non-interactive shells. `$PWD` is set by
// sh to the CWD OpenCode chose when spawning the MCP child. `exec`
// replaces the shell so the MCP child becomes ctx itself, with no
// sh process layered between OpenCode and the JSON-RPC stream.
//
// Returns:
//   - []string: argv suitable for OpenCode's McpLocalConfig.command.
func launchCommand() []string {
	bin := mcpServer.Command
	if resolved, err := exec.LookPath(bin); err == nil {
		if abs, absErr := filepath.Abs(resolved); absErr == nil {
			bin = abs
		}
	}
	binAndArgs := append([]string{bin}, mcpServer.Args()...)
	script := fmt.Sprintf(
		cfgShell.FormatPOSIXSpawnRelativeCtxDir,
		env.CtxDir, cfgDir.Context,
		strings.Join(binAndArgs, token.Space),
	)
	return []string{cfgShell.Sh, cfgShell.CmdFlag, script}
}

// globalConfigPath returns the path to the OpenCode global config
// file (~/.config/opencode/opencode.json or $OPENCODE_HOME/opencode.json).
//
// Returns:
//   - string: absolute path to the OpenCode config file.
//   - error: non-nil when the user home directory cannot be resolved.
func globalConfigPath() (string, error) {
	ocHome := os.Getenv(cfgHook.EnvOpenCodeHome)
	if ocHome == "" {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return "", homeErr
		}
		ocHome = filepath.Join(
			home, cfgHook.DirXDGConfig, cfgHook.DirOpenCodeHome,
		)
	}
	return filepath.Join(ocHome, cfgHook.FileOpenCodeJSON), nil
}

// ensureMCPConfig registers the ctx MCP server in the OpenCode
// global config file (~/.config/opencode/opencode.json).
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
	target, pathErr := globalConfigPath()
	if pathErr != nil {
		return pathErr
	}

	existing := make(map[string]interface{})
	data, readErr := ctxIo.SafeReadUserFile(target)
	if readErr == nil && len(bytes.TrimSpace(data)) > 0 {
		if jErr := json.Unmarshal(data, &existing); jErr != nil {
			return jErr
		}
	}

	servers, _ := existing[cfgHook.KeyMCP].(map[string]interface{})
	if servers == nil {
		servers = make(map[string]interface{})
	}

	if _, ok := servers[mcpServer.Name]; ok {
		writeSetup.InfoOpenCodeSkipped(cmd, target)
		return nil
	}

	servers[mcpServer.Name] = map[string]interface{}{
		cfgHook.KeyType:    cfgHook.MCPServerType,
		cfgHook.KeyCommand: launchCommand(),
		cfgHook.KeyEnabled: true,
	}
	existing[cfgHook.KeyMCP] = servers

	// Ensure the directory exists.
	dir := filepath.Dir(target)
	if mkErr := ctxIo.SafeMkdirAll(dir, fs.PermExec); mkErr != nil {
		return mkErr
	}

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
