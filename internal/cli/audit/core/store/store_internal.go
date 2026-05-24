//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package store

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/cli/audit/core/parse"
	cfgAudit "github.com/ActiveMemory/ctx/internal/config/audit"
	errAudit "github.com/ActiveMemory/ctx/internal/err/audit"
	"github.com/ActiveMemory/ctx/internal/io"
)

// loadOne reads and parses a single report by id from a
// pre-resolved audit directory.
//
// Parameters:
//   - dir: absolute audit-directory path
//   - id: report basename without extension
//
// Returns:
//   - Report: parsed report
//   - error: [errAudit.UnknownID] on absence, or any
//     read / parse error
func loadOne(dir, id string) (Report, error) {
	path := filepath.Join(dir, id+cfgAudit.ReportExt)
	data, readErr := io.SafeReadFile(dir, id+cfgAudit.ReportExt)
	if errors.Is(readErr, os.ErrNotExist) {
		return Report{}, errAudit.UnknownID(id)
	}
	if readErr != nil {
		return Report{}, errAudit.ReadReport(id, readErr)
	}

	header, body, parseErr := parse.Frontmatter(data)
	if parseErr != nil {
		return Report{}, errAudit.ParseReport(id, parseErr)
	}

	return Report{
		ID:          id,
		Path:        path,
		Kind:        header.Kind,
		Status:      header.Status,
		CommitRange: header.CommitRange,
		GeneratedAt: header.GeneratedAt,
		Generator:   header.Generator,
		Digest:      header.Digest,
		Body:        body,
	}, nil
}
