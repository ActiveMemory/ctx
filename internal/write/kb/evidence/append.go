//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package evidence

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgKbEvidence "github.com/ActiveMemory/ctx/internal/config/kb/evidence"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errKbEvidence "github.com/ActiveMemory/ctx/internal/err/kb/evidence"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Append writes one EV row to the evidence-index file at path.
// When row.ID is empty, the file's high-water mark is consulted
// and the next sequential ID is allocated.
//
// Parameters:
//   - path: full path to `.context/kb/evidence-index.md`.
//   - row: row to append. ID is populated on return when
//     allocated.
//
// Returns:
//   - Row: the appended row with ID populated.
//   - error: [errKbEvidence.ErrDuplicateID],
//     [errKbEvidence.ErrInvalidBand], or wrapped I/O errors.
func Append(path string, row Row) (Row, error) {
	if bandErr := validateBand(row.Confidence); bandErr != nil {
		return Row{}, bandErr
	}

	existing, readErr := ctxIo.SafeReadUserFile(path)
	if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
		return Row{}, errKbEvidence.ReadIndex(readErr)
	}

	if row.ID == "" {
		next, scanErr := maxIDFrom(string(existing))
		if scanErr != nil {
			return Row{}, scanErr
		}
		row.ID = formatID(next + 1)
	} else if alreadyExists(string(existing), row.ID) {
		return Row{}, errKbEvidence.DuplicateID(row.ID)
	}
	if row.Extracted.IsZero() {
		row.Extracted = time.Now().UTC().Truncate(time.Second)
	}

	if mkErr := ctxIo.SafeMkdirAll(
		filepath.Dir(path), cfgFs.PermExec,
	); mkErr != nil {
		return Row{}, errKbEvidence.MkdirDir(mkErr)
	}

	needsHeader := len(existing) == 0
	f, openErr := ctxIo.SafeAppendFile(path, cfgFs.PermSecret)
	if openErr != nil {
		return Row{}, errKbEvidence.OpenIndex(openErr)
	}
	defer func() { _ = f.Close() }()

	var sb strings.Builder
	if needsHeader {
		sb.WriteString(cfgKbEvidence.TitleHeading)
		sb.WriteString(token.DoubleNewline)
		sb.WriteString(cfgKbEvidence.LeadParagraph1)
		sb.WriteString(token.NewlineLF)
		sb.WriteString(cfgKbEvidence.LeadParagraph2)
		sb.WriteString(token.DoubleNewline)
		sb.WriteString(cfgKbEvidence.TableHeader)
		sb.WriteString(token.NewlineLF)
	}
	sb.WriteString(renderRow(row))
	sb.WriteString(token.NewlineLF)
	if _, writeErr := f.WriteString(sb.String()); writeErr != nil {
		return Row{}, errKbEvidence.WriteRow(writeErr)
	}
	return row, nil
}
