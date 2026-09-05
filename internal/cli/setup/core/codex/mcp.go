//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/codex"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// ensureMCPConfig registers the ctx MCP server in the project
// .codex/config.toml. Creates the file when absent, appends the
// [mcp_servers.ctx] table when its header is missing, and skips
// when the header is present. Existing bytes are never rewritten.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if the target is not a regular file or the
//     write fails
func ensureMCPConfig(cmd *cobra.Command) error {
	target := cfgSetup.MCPConfigPathCodex
	if _, validateErr := validateManagedTarget(target); validateErr != nil {
		return validateErr
	}

	existing, readErr := ctxIo.SafeReadUserFile(target)
	if readErr != nil && !os.IsNotExist(readErr) {
		return errFs.FileRead(target, readErr)
	}

	out, outcome := codex.EnsureMCPTable(existing)
	if outcome == codex.OutcomeSkipped {
		writeSetup.InfoCodexSkipped(cmd, target)
		return nil
	}

	if writeFileErr := writeManaged(target, out); writeFileErr != nil {
		return writeFileErr
	}
	if outcome == codex.OutcomeMerged {
		writeSetup.InfoCodexMerged(cmd, target)
		return nil
	}
	writeSetup.InfoCodexCreated(cmd, target)
	return nil
}
