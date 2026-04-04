//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package kiro

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	mcpServer "github.com/ActiveMemory/ctx/internal/config/mcp/server"
	cfgMcpTool "github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errSetup "github.com/ActiveMemory/ctx/internal/err/setup"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/steering"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// ensureMCPConfig creates .kiro/settings/mcp.json with
// the ctx MCP server config. Skips if the file exists.
func ensureMCPConfig(cmd *cobra.Command) error {
	settingsDir := filepath.Join(
		cfgSetup.DirKiro, cfgSetup.DirSettings,
	)
	target := filepath.Join(settingsDir, cfgSetup.FileMCPJSON)

	if _, statErr := ctxIo.SafeStat(
		target,
	); statErr == nil {
		writeSetup.DeployFileExists(cmd, target)
		return nil
	}

	if mkdirErr := ctxIo.SafeMkdirAll(
		settingsDir, fs.PermExec,
	); mkdirErr != nil {
		return errSetup.CreateDir(settingsDir, mkdirErr)
	}

	cfg := mcpConfig{
		MCPServers: map[string]serverEntry{
			mcpServer.Name: {
				Command:  mcpServer.Command,
				Args:     mcpServer.Args(),
				Disabled: false,
				AutoApprove: []string{
					cfgMcpTool.Status,
					cfgMcpTool.SteeringGet,
					cfgMcpTool.Search,
					cfgMcpTool.SessionStart,
					cfgMcpTool.SessionEnd,
					cfgMcpTool.Next,
					cfgMcpTool.Remind,
				},
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

// syncSteering syncs steering files to Kiro format
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
		steeringDir, token.Dot, cfgHook.ToolKiro,
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
