//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package shared

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/hub"
	"github.com/ActiveMemory/ctx/internal/io"
	writeServe "github.com/ActiveMemory/ctx/internal/write/serve"
)

// defaultPort is the default hub listen port.
const defaultPort = 9900

// hubDataDir is the subdirectory for hub data files.
const hubDataDir = "hub-data"

// adminTokenFile stores the admin token after first run.
const adminTokenFile = "admin.token"

// dataDirPerm is the permission for the hub data directory.
const dataDirPerm = fs.PermKeyDir

// hubDir returns the default hub data directory path.
// Uses ~/.ctx/hub-data/ (same parent as the encryption key).
func hubDir() (string, error) {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", homeErr
	}
	p := filepath.Join(home, dir.CtxData, hubDataDir)
	return p, io.SafeMkdirAll(p, fs.PermKeyDir)
}

// loadOrCreateAdmin loads an existing admin token or
// generates a new one on first run.
func loadOrCreateAdmin(
	cmd *cobra.Command, dataDir string,
) (string, error) {
	tokenPath := filepath.Join(dataDir, adminTokenFile)
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
