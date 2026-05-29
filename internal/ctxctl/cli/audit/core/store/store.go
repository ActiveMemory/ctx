//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	cfgAudit "github.com/ActiveMemory/ctx/internal/ctxctl/config/audit"
	errAudit "github.com/ActiveMemory/ctx/internal/ctxctl/err/audit"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Dir returns the absolute path of the audit directory
// (`.context/<DirName>/`).
//
// Returns:
//   - string: audit directory path
//   - error: propagated from [rc.ContextDir]
func Dir() (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ctxDir, cfgAudit.DirName), nil
}

// DismissedLedgerPath returns the absolute path of the
// dismissal-ledger JSON file.
//
// Returns:
//   - string: ledger path
//   - error: propagated from [Dir]
func DismissedLedgerPath() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, cfgAudit.DismissedLedger), nil
}

// Read loads every audit report from disk in id order.
//
// Returns:
//   - []Report: parsed reports (nil when the audit
//     directory does not exist yet)
//   - error: non-nil on directory-read, per-file read, or
//     frontmatter parse failure
func Read() ([]Report, error) {
	dir, dirErr := Dir()
	if dirErr != nil {
		return nil, dirErr
	}

	entries, readErr := os.ReadDir(dir) //nolint:gosec // dir is rc-derived
	if errors.Is(readErr, os.ErrNotExist) {
		return nil, nil
	}
	if readErr != nil {
		return nil, errAudit.ReadReport(cfgAudit.DirName, readErr)
	}

	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, cfgAudit.ReportExt) {
			continue
		}
		if strings.HasPrefix(name, token.Dot) {
			continue
		}
		ids = append(ids, strings.TrimSuffix(name, cfgAudit.ReportExt))
	}
	sort.Strings(ids)

	out := make([]Report, 0, len(ids))
	for _, id := range ids {
		r, loadErr := loadOne(dir, id)
		if loadErr != nil {
			return nil, loadErr
		}
		out = append(out, r)
	}
	return out, nil
}

// ReadOne loads a single report by id.
//
// Parameters:
//   - id: report basename without extension
//
// Returns:
//   - Report: parsed report
//   - error: [errAudit.UnknownID] when missing, or any
//     read / parse error
func ReadOne(id string) (Report, error) {
	dir, dirErr := Dir()
	if dirErr != nil {
		return Report{}, dirErr
	}
	return loadOne(dir, id)
}

// ReadDismissals loads the dismissal ledger.
//
// Returns:
//   - DismissalLedger: parsed ledger (empty entries map
//     when the file is missing)
//   - error: non-nil on read or unmarshal failure
func ReadDismissals() (DismissalLedger, error) {
	path, pathErr := DismissedLedgerPath()
	if pathErr != nil {
		return DismissalLedger{}, pathErr
	}
	data, readErr := io.SafeReadUserFile(path)
	if errors.Is(readErr, os.ErrNotExist) {
		return DismissalLedger{
			Entries: map[string]DismissedAt{},
		}, nil
	}
	if readErr != nil {
		return DismissalLedger{}, errAudit.ReadDismissal(readErr)
	}

	var led DismissalLedger
	if jsonErr := json.Unmarshal(data, &led); jsonErr != nil {
		return DismissalLedger{}, errAudit.ReadDismissal(jsonErr)
	}
	if led.Entries == nil {
		led.Entries = map[string]DismissedAt{}
	}
	return led, nil
}

// WriteDismissals persists the dismissal ledger atomically.
//
// Parameters:
//   - led: ledger to persist
//
// Returns:
//   - error: non-nil on marshal or write failure
func WriteDismissals(led DismissalLedger) error {
	path, pathErr := DismissedLedgerPath()
	if pathErr != nil {
		return pathErr
	}

	if mkErr := io.SafeMkdirAll(
		filepath.Dir(path), fs.PermRestrictedDir,
	); mkErr != nil {
		return errAudit.WriteDismissal(mkErr)
	}

	data, marshalErr := json.MarshalIndent(
		led, "", token.Indent2,
	)
	if marshalErr != nil {
		return errAudit.WriteDismissal(marshalErr)
	}
	if writeErr := io.SafeWriteFile(
		path, data, fs.PermFile,
	); writeErr != nil {
		return errAudit.WriteDismissal(writeErr)
	}
	return nil
}

// Dismiss marks one report id as dismissed against its
// current digest. Returns [errAudit.UnknownID] if the report
// does not exist.
//
// Parameters:
//   - id: report basename
//
// Returns:
//   - error: non-nil on missing report or ledger I/O
func Dismiss(id string) error {
	r, readErr := ReadOne(id)
	if readErr != nil {
		return readErr
	}
	led, ledErr := ReadDismissals()
	if ledErr != nil {
		return ledErr
	}
	led.Entries[id] = DismissedAt{
		Digest: r.Digest,
		At:     time.Now().UTC(),
	}
	return WriteDismissals(led)
}

// DismissAll marks every current report as dismissed.
// Returns the count actually dismissed.
//
// Returns:
//   - int: number of reports dismissed
//   - error: non-nil on read or ledger write failure
func DismissAll() (int, error) {
	reports, readErr := Read()
	if readErr != nil {
		return 0, readErr
	}
	if len(reports) == 0 {
		return 0, nil
	}
	led, ledErr := ReadDismissals()
	if ledErr != nil {
		return 0, ledErr
	}
	now := time.Now().UTC()
	for _, r := range reports {
		led.Entries[r.ID] = DismissedAt{
			Digest: r.Digest,
			At:     now,
		}
	}
	if writeErr := WriteDismissals(led); writeErr != nil {
		return 0, writeErr
	}
	return len(reports), nil
}

// IsDismissed reports whether the given report is dismissed
// against its current digest. A dismissal recorded against
// an older digest does not count: a fresh audit run with
// new findings re-surfaces the report.
//
// Parameters:
//   - r: candidate report
//   - led: dismissal ledger
//
// Returns:
//   - bool: true when dismissed against the current digest
func IsDismissed(r Report, led DismissalLedger) bool {
	entry, ok := led.Entries[r.ID]
	if !ok {
		return false
	}
	return entry.Digest == r.Digest
}
