//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	mcpServer "github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errSetup "github.com/ActiveMemory/ctx/internal/err/setup"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// launchCommand returns the OpenCode `command` array for the ctx
// MCP server. The emitted argv is:
//
//	[/abs/path/to/ctx, mcp, serve]
//
// Under the cwd-anchored resolution model
// (spec: specs/cwd-anchored-context.md) ctx anchors to its
// working directory, so the OpenCode-chosen CWD is sufficient
// project context; no env-var wrapper or shell layer is needed.
// The binary path is resolved to an absolute path at setup time
// via exec.LookPath so OpenCode can spawn the MCP child regardless
// of the PATH inherited by non-interactive shells.
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
	return append([]string{bin}, mcpServer.Args()...)
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
// Merge-safe: reads existing config, adds or updates the ctx
// server under the "mcp" key, writes back, and preserves all
// unrelated config keys. Treats a missing or empty file as "no
// existing config" rather than an error.
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

	if _, validateErr := validateManagedTarget(target); validateErr != nil {
		return validateErr
	}

	existing := make(map[string]interface{})
	data, readErr := ctxIo.SafeReadUserFile(target)
	if readErr == nil {
		if len(bytes.TrimSpace(data)) > 0 {
			if jErr := json.Unmarshal(data, &existing); jErr != nil {
				return errSetup.MarshalConfig(jErr)
			}
		}
	} else if !os.IsNotExist(readErr) {
		return errFs.FileRead(target, readErr)
	}

	servers, _ := existing[cfgHook.KeyMCP].(map[string]interface{})
	if servers == nil {
		servers = make(map[string]interface{})
	}

	newServer := map[string]interface{}{
		cfgHook.KeyType:    cfgHook.MCPServerType,
		cfgHook.KeyCommand: launchCommand(),
		cfgHook.KeyEnabled: true,
	}
	if existingServer, ok := servers[mcpServer.Name]; ok {
		if existingMap, mapOK := existingServer.(map[string]interface{}); mapOK {
			current, marshalErr := json.Marshal(existingMap)
			if marshalErr != nil {
				return errSetup.MarshalConfig(marshalErr)
			}
			expected, marshalErr := json.Marshal(newServer)
			if marshalErr != nil {
				return errSetup.MarshalConfig(marshalErr)
			}
			if bytes.Equal(current, expected) {
				writeSetup.InfoOpenCodeSkipped(cmd, target)
				return nil
			}
		}
	}

	servers[mcpServer.Name] = newServer
	existing[cfgHook.KeyMCP] = servers

	// Ensure the directory exists.
	dir := filepath.Dir(target)
	if mkErr := ctxIo.SafeMkdirAll(dir, fs.PermExec); mkErr != nil {
		return errFs.Mkdir(dir, mkErr)
	}

	out, marshalErr := json.MarshalIndent(
		existing, "", token.Indent2,
	)
	if marshalErr != nil {
		return errSetup.MarshalConfig(marshalErr)
	}
	out = append(out, token.NewlineLF...)

	if writeFileErr := ctxIo.SafeWriteFileAtomic(
		target, out, fs.PermFile,
	); writeFileErr != nil {
		return errFs.FileWrite(target, writeFileErr)
	}
	writeSetup.InfoOpenCodeCreated(cmd, target)
	return nil
}
