//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"github.com/ActiveMemory/ctx/internal/hub"
	"github.com/ActiveMemory/ctx/internal/io"
	writeServe "github.com/ActiveMemory/ctx/internal/write/serve"
)

// defaultPort is the default hub listen port.
const defaultPort = 9900

// resolveDataDir returns the data directory, creating it
// if needed.
//
// Parameters:
//   - dataDir: Explicit data dir path, or empty for default
//
// Returns:
//   - string: Resolved absolute data directory path
//   - error: Non-nil on mkdir failure
func resolveDataDir(dataDir string) (string, error) {
	if dataDir == "" {
		return defaultDataDir()
	}
	return dataDir, io.SafeMkdirAll(
		dataDir, fs.PermKeyDir,
	)
}

// defaultDataDir returns the default hub data directory path.
// Uses ~/.ctx/hub-data/ (same parent as the encryption key).
//
// Returns:
//   - string: Absolute path to ~/.ctx/hub-data/
//   - error: Non-nil on home-dir lookup or mkdir failure
func defaultDataDir() (string, error) {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", homeErr
	}
	p := filepath.Join(home, dir.CtxData, cfgHub.DirHubData)
	return p, io.SafeMkdirAll(p, fs.PermKeyDir)
}

// loadOrCreateAdmin loads an existing admin token or
// generates a new one on first run.
//
// Parameters:
//   - cmd: Cobra command for output (prints token on first run)
//   - dataDir: Hub data directory containing admin.token
//
// Returns:
//   - string: Admin token (existing or newly generated)
//   - error: Non-nil on generation or I/O failure
func loadOrCreateAdmin(
	cmd *cobra.Command, dataDir string,
) (string, error) {
	tokenPath := filepath.Join(dataDir, cfgHub.FileAdminToken)
	data, readErr := io.SafeReadUserFile(tokenPath)
	if readErr == nil && len(data) > 0 {
		return string(data), nil
	}

	adminToken, genErr := hub.GenerateAdminToken()
	if genErr != nil {
		return "", genErr
	}

	if writeErr := io.SafeWriteFile(
		tokenPath, []byte(adminToken), fs.PermSecret,
	); writeErr != nil {
		return "", writeErr
	}

	writeServe.AdminToken(cmd, adminToken)

	return adminToken, nil
}
