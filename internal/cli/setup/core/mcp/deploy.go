//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errSetup "github.com/ActiveMemory/ctx/internal/err/setup"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/steering"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Deploy writes an MCP config file if it does not already
// exist. It creates the parent directory, marshals the
// config as indented JSON, and prints a confirmation
// message.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dir: Parent directory for the config file
//   - filename: Config file name (e.g. "mcp.json")
//   - cfg: Config struct to marshal as JSON
//
// Returns:
//   - error: Non-nil on directory creation, marshal,
//     or write failure
func Deploy(
	cmd *cobra.Command,
	dir, filename string,
	cfg any,
) error {
	target := filepath.Join(dir, filename)

	if _, statErr := ctxIo.SafeStat(
		target,
	); statErr == nil {
		writeSetup.DeployFileExists(cmd, target)
		return nil
	}

	if mkdirErr := ctxIo.SafeMkdirAll(
		dir, fs.PermExec,
	); mkdirErr != nil {
		return errSetup.CreateDir(dir, mkdirErr)
	}

	data, marshalErr := json.MarshalIndent(
		cfg, "", token.Indent2,
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

// SyncSteering syncs steering files to a tool-native
// format. If no steering directory exists, prints a
// message and returns nil.
//
// Parameters:
//   - cmd: Cobra command for output
//   - tool: Tool identifier for sync format selection
//
// Returns:
//   - error: Non-nil on sync failure
func SyncSteering(
	cmd *cobra.Command, tool string,
) error {
	steeringDir := rc.SteeringDir()
	if _, statErr := ctxIo.SafeStat(
		steeringDir,
	); os.IsNotExist(statErr) {
		writeSetup.DeployNoSteering(cmd)
		return nil
	}

	report, syncErr := steering.SyncTool(
		steeringDir, token.Dot, tool,
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
