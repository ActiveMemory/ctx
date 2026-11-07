//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agy

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/agy"
	cfgAsset "github.com/ActiveMemory/ctx/internal/config/asset"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// ensureMCPConfig registers the ctx MCP server in the Antigravity
// mcp_config.json configuration file (<ConfigDir()>/mcp_config.json).
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: non-nil if reading, merging, or writing fails
func ensureMCPConfig(cmd *cobra.Command) error {
	configDir := agy.ConfigDir()
	if configDir == "" {
		_, homeErr := os.UserHomeDir()
		return errInit.HomeDir(homeErr)
	}
	target := filepath.Join(configDir, cfgAsset.FileMCPConfigJSON)
	if _, validateErr := validateManagedTarget(target); validateErr != nil {
		return validateErr
	}

	existing, readErr := ctxIo.SafeReadUserFile(target)
	if readErr != nil && !os.IsNotExist(readErr) {
		return errFs.FileRead(target, readErr)
	}

	out, outcome, mergeErr := agy.EnsureMCPConfig(existing)
	if mergeErr != nil {
		return mergeErr
	}
	if outcome == agy.OutcomeSkipped {
		writeSetup.InfoAgySkipped(cmd, target)
		return nil
	}

	if writeFileErr := writeManaged(target, out); writeFileErr != nil {
		return writeFileErr
	}
	if outcome == agy.OutcomeMerged {
		writeSetup.InfoAgyMerged(cmd, target)
		return nil
	}
	writeSetup.InfoAgyCreated(cmd, target)
	return nil
}
