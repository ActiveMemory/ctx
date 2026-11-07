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
	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeErr "github.com/ActiveMemory/ctx/internal/write/err"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// deployHooks creates or merges .agents/hooks.json from the embedded
// Antigravity hooks manifest.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - bool: true when written or up to date, false if unparseable
//   - error: non-nil if validation or write fails
func deployHooks(cmd *cobra.Command) (bool, error) {
	target := cfgSetup.HooksPathAgy
	if _, validateErr := validateManagedTarget(target); validateErr != nil {
		return false, validateErr
	}

	embedded, assetErr := agent.AgyHooksJSON()
	if assetErr != nil {
		return false, assetErr
	}

	existing, readErr := ctxIo.SafeReadUserFile(target)
	if readErr != nil && !os.IsNotExist(readErr) {
		return false, errFs.FileRead(target, readErr)
	}

	out, outcome, mergeErr := agy.MergeHooks(existing, embedded)
	if mergeErr != nil {
		writeErr.WarnFile(cmd, target, mergeErr)
		return false, nil
	}
	if outcome == agy.OutcomeSkipped {
		writeSetup.InfoAgySkipped(cmd, target)
		return true, nil
	}

	if writeFileErr := writeManaged(target, out); writeFileErr != nil {
		return false, writeFileErr
	}
	if outcome == agy.OutcomeMerged {
		writeSetup.InfoAgyMerged(cmd, target)
		return true, nil
	}
	writeSetup.InfoAgyCreated(cmd, target)
	return true, nil
}

// writeManaged writes a managed file, creating its parent directory
// when needed.
//
// Parameters:
//   - target: file path
//   - content: bytes to write
//
// Returns:
//   - error: Non-nil if directory creation or write fails
func writeManaged(target string, content []byte) error {
	dir := filepath.Dir(target)
	if mkErr := ctxIo.SafeMkdirAll(dir, fs.PermExec); mkErr != nil {
		return errFs.Mkdir(dir, mkErr)
	}
	if wErr := ctxIo.SafeWriteFileAtomic(
		target, content, fs.PermFile,
	); wErr != nil {
		return errFs.FileWrite(target, wErr)
	}
	return nil
}
