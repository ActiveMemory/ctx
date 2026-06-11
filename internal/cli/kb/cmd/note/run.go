//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package note

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	kbPath "github.com/ActiveMemory/ctx/internal/cli/kb/core/path"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	errKbCli "github.com/ActiveMemory/ctx/internal/err/kb/cli"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Run appends a timestamped one-liner to
// `.context/ingest/findings.md`. Refuses on empty input.
//
// Parameters:
//   - cobraCmd: cobra command for output.
//   - noteText: free-text note body (trimmed; empty rejected).
//
// Returns:
//   - error: refusal on empty input, or wrapped I/O failure.
func Run(cobraCmd *cobra.Command, noteText string) (err error) {
	if noteText == "" {
		cobraCmd.SilenceUsage = true
		return errKbCli.ErrNoteNoText
	}
	ingestDir, dirErr := kbPath.IngestDir()
	if dirErr != nil {
		return dirErr
	}
	findings := filepath.Join(ingestDir, cfgKB.Findings)
	mkErr := io.SafeMkdirAll(ingestDir, cfgFs.PermExec)
	if mkErr != nil {
		return errKbCli.MkdirIngest(mkErr)
	}
	f, openErr := io.SafeAppendFile(findings, cfgFs.PermSecret)
	if openErr != nil {
		return errKbCli.OpenFindings(openErr)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = errKbCli.WriteFinding(cerr)
		}
	}()

	stamp := time.Now().UTC().Format(time.RFC3339)
	line := fmt.Sprintf(
		desc.Text(text.DescKeyWriteKbFindingLine), stamp, noteText,
	)
	if _, writeErr := f.WriteString(line); writeErr != nil {
		return errKbCli.WriteFinding(writeErr)
	}
	io.SafeFprintf(
		cobraCmd.OutOrStdout(),
		desc.Text(text.DescKeyWriteKbAppendedTo),
		findings,
	)
	return nil
}
