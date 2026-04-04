//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cline

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	mcpServer "github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/config/token"
	cfgVscode "github.com/ActiveMemory/ctx/internal/config/vscode"
	errSetup "github.com/ActiveMemory/ctx/internal/err/setup"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/steering"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// ensureMCPConfig creates .vscode/mcp.json with the ctx
// MCP server configuration. Skips if the file exists.
func ensureMCPConfig(cmd *cobra.Command) error {
	target := filepath.Join(
		cfgVscode.Dir, cfgVscode.FileMCPJSON,
	)

	if _, statErr := ctxIo.SafeStat(target); statErr == nil {
		writeSetup.DeployFileExists(cmd, target)
		return nil
	}

	if mkdirErr := ctxIo.SafeMkdirAll(
		cfgVscode.Dir, fs.PermExec,
	); mkdirErr != nil {
		return errSetup.CreateDir(cfgVscode.Dir, mkdirErr)
	}

	cfg := vscodeMCPConfig{
		Servers: map[string]vscodeMCPServer{
			mcpServer.Name: {
				Command: mcpServer.Command,
				Args:    mcpServer.Args(),
			},
		},
	}

	data, marshalErr := json.MarshalIndent(
		cfg, "", "  ",
	)
	if marshalErr != nil {
		return errSetup.MarshalConfig(marshalErr)
	}
	data = append(data, token.NewlineLF[0])

	if writeErr := ctxIo.SafeWriteFile(
		target, data, fs.PermFile,
	); writeErr != nil {
		return errSetup.WriteFile(target, writeErr)
	}

	writeSetup.DeployFileCreated(cmd, target)
	return nil
}

// syncSteering syncs steering files to Cline format
// if a steering directory exists.
func syncSteering(cmd *cobra.Command) error {
	steeringDir := rc.SteeringDir()
	if _, statErr := ctxIo.SafeStat(
		steeringDir,
	); os.IsNotExist(statErr) {
		writeSetup.DeployNoSteering(cmd)
		return nil
	}

	report, syncErr := steering.SyncTool(
		steeringDir, token.Dot, cfgHook.ToolCline,
	)
	if syncErr != nil {
		return errSetup.SyncSteering(syncErr)
	}

	for _, name := range report.Written {
		writeSetup.DeploySteeringSynced(cmd, name)
	}
	for _, name := range report.Skipped {
		writeSetup.DeploySteeringSkipped(cmd, name)
	}
	return nil
}
