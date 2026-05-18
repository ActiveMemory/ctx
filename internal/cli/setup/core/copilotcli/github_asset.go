//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package copilotcli

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// deployGithubAsset writes a single Copilot CLI markdown asset
// under `.github/<subDir>/<fileName>`, idempotently. If the
// target already exists it logs a skip notice and returns nil;
// otherwise it creates the directory, fetches the embedded
// content via readContent, writes the file, and logs a created
// notice.
//
// Parameters:
//   - cmd: cobra command for output messages.
//   - subDir: subdirectory under `.github/` (e.g.
//     [cfgHook.DirGitHubAgents]).
//   - fileName: target filename inside subDir.
//   - readContent: embedded-asset accessor that returns the
//     file body.
//
// Returns:
//   - error: non-nil on directory creation, embedded read, or
//     file write failure.
func deployGithubAsset(
	cmd *cobra.Command,
	subDir, fileName string,
	readContent func() ([]byte, error),
) error {
	assetDir := filepath.Join(cfgHook.DirGitHub, subDir)
	target := filepath.Join(assetDir, fileName)

	if _, statErr := os.Stat(target); statErr == nil {
		writeSetup.InfoCopilotCLISkipped(cmd, target)
		return nil
	}

	if mkErr := ctxIo.SafeMkdirAll(assetDir, fs.PermExec); mkErr != nil {
		return mkErr
	}

	content, readErr := readContent()
	if readErr != nil {
		return readErr
	}
	if wErr := ctxIo.SafeWriteFile(target, content, fs.PermFile); wErr != nil {
		return wErr
	}
	writeSetup.InfoCopilotCLICreated(cmd, target)
	return nil
}
