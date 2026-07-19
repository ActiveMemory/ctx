//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package apply

import (
	"encoding/json"
	"io"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/disclosure"
	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	writeDisc "github.com/ActiveMemory/ctx/internal/write/disclosure"
)

// Run executes the apply command: resolve the file's kind, read the
// digest plan (from --plan or stdin), move its staged entries into theme
// files, fold their gists into the root, and report what moved. The move
// itself is guarded by internal/disclosure.Apply (append→verify→remove);
// on any failure the root is left byte-identical.
//
// Parameters:
//   - cmd: Cobra command for the I/O streams
//   - path: path to the canonical knowledge file
//   - planPath: digest plan path, or "-" for stdin
//   - jsonOutput: if true, emit the result as JSON
//
// Returns:
//   - error: NotAKnowledgeFile, a read error, a JSON error, a disclosure
//     guard sentinel, or nil
func Run(cmd *cobra.Command, path, planPath string, jsonOutput bool) error {
	clean := filepath.Clean(path)
	if _, ok := disclosure.KindFor(filepath.Base(clean)); !ok {
		cmd.SilenceUsage = true
		return errDisc.NotAKnowledgeFile(path)
	}

	var planBytes []byte
	if planPath == token.Dash {
		b, stdinErr := io.ReadAll(cmd.InOrStdin())
		if stdinErr != nil {
			cmd.SilenceUsage = true
			return errFs.FileRead(token.Dash, stdinErr)
		}
		planBytes = b
	} else {
		b, planReadErr := internalIo.SafeReadUserFile(filepath.Clean(planPath))
		if planReadErr != nil {
			cmd.SilenceUsage = true
			return errFs.FileRead(planPath, planReadErr)
		}
		planBytes = b
	}

	var plan disclosure.Plan
	if jsonErr := json.Unmarshal(planBytes, &plan); jsonErr != nil {
		cmd.SilenceUsage = true
		return jsonErr
	}

	res, applyErr := disclosure.Apply(clean, plan, filepath.Dir(clean))
	if applyErr != nil {
		cmd.SilenceUsage = true
		return applyErr
	}

	if jsonOutput {
		return writeDisc.ApplyJSON(cmd, res)
	}
	writeDisc.ApplyHuman(cmd, res)
	return nil
}
