//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"encoding/json"
	"os"
	"strings"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgToken "github.com/ActiveMemory/ctx/internal/config/token"
	errDream "github.com/ActiveMemory/ctx/internal/err/dream"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// AppendLedger appends one disposition to the append-only ledger at
// <dreamsDir>/ledger.md, creating the notebook directory if needed.
// Each entry is written as a Markdown list item whose payload is the
// JSON encoding of the LedgerEntry — human-readable as a list, and
// machine-parseable by ReadLedger. The ledger is never rewritten, only
// appended, so the decision trail is tamper-evident.
//
// Parameters:
//   - dreamsDir: absolute path to the dreams/ notebook directory
//   - entry: the disposition to record
//
// Returns:
//   - error: non-nil on directory creation, JSON marshal, or append
//     failure
func AppendLedger(dreamsDir string, entry LedgerEntry) error {
	if mkErr := ctxIo.SafeMkdirAll(
		dreamsDir, cfgFs.PermRestrictedDir,
	); mkErr != nil {
		return errDream.Mkdir(dreamsDir, mkErr)
	}
	payload, marshalErr := json.Marshal(entry)
	if marshalErr != nil {
		return errDream.MarshalEntry(marshalErr)
	}
	line := cfgToken.PrefixListDash + string(payload) + cfgToken.NewlineLF
	path := ledgerPath(dreamsDir)
	if appendErr := ctxIo.AppendBytes(
		path, []byte(line), cfgFs.PermSecret,
	); appendErr != nil {
		return errDream.AppendLedger(path, appendErr)
	}
	return nil
}

// ReadLedger reads back every disposition recorded in the ledger at
// <dreamsDir>/ledger.md, in append order. A missing ledger is not an
// error: it yields an empty slice. Lines that are not JSON list-item
// payloads are skipped, so prose interleaved into the notebook does not
// corrupt the readback.
//
// Parameters:
//   - dreamsDir: absolute path to the dreams/ notebook directory
//
// Returns:
//   - []LedgerEntry: the recorded dispositions in append order
//   - error: non-nil on a read failure other than not-exist
func ReadLedger(dreamsDir string) ([]LedgerEntry, error) {
	path := ledgerPath(dreamsDir)
	data, readErr := ctxIo.SafeReadUserFile(path)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return []LedgerEntry{}, nil
		}
		return nil, errDream.ReadLedger(path, readErr)
	}
	var entries []LedgerEntry
	for _, raw := range strings.Split(string(data), cfgToken.NewlineLF) {
		line := strings.TrimSpace(raw)
		if !strings.HasPrefix(line, cfgToken.PrefixListDash) {
			continue
		}
		payload := strings.TrimPrefix(line, cfgToken.PrefixListDash)
		var entry LedgerEntry
		if json.Unmarshal([]byte(payload), &entry) != nil {
			continue
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// Seen is the dedup-against-seen signal: it reports whether the ledger
// already records a disposition for proposalID. A proposal whose source
// has not changed and that has already been decided (accepted, rejected,
// amended, or skipped) is not re-surfaced. Rejections count as seen, by
// design — the dream does not re-nag a rejected disposition unless the
// source content changes.
//
// Parameters:
//   - entries: ledger entries (from ReadLedger)
//   - proposalID: the proposal ID to test
//
// Returns:
//   - bool: true when a disposition for proposalID exists in entries
func Seen(entries []LedgerEntry, proposalID string) bool {
	for _, e := range entries {
		if e.ProposalID == proposalID {
			return true
		}
	}
	return false
}
