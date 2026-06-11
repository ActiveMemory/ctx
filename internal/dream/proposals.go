//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	errDream "github.com/ActiveMemory/ctx/internal/err/dream"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// LatestRunDir returns the most recent per-run directory under
// dreamsDir (the lexically greatest timestamped name, since the run
// layout is sortable). A dreams/ directory with no run subdirectories
// yields an empty string and a nil error.
//
// Parameters:
//   - dreamsDir: absolute path to the dreams/ notebook directory
//
// Returns:
//   - string: absolute path to the latest run directory, or "" when
//     none exist
//   - error: non-nil on a read failure other than not-exist
func LatestRunDir(dreamsDir string) (string, error) {
	entries, readErr := os.ReadDir(dreamsDir)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return "", nil
		}
		return "", errDream.ReadProposals(dreamsDir, readErr)
	}
	var runs []string
	for _, e := range entries {
		if e.IsDir() && runDirName(e.Name()) {
			runs = append(runs, e.Name())
		}
	}
	if len(runs) == 0 {
		return "", nil
	}
	sort.Strings(runs)
	return filepath.Join(dreamsDir, runs[len(runs)-1]), nil
}

// ReadProposals reads and decodes the proposals file the executor
// wrote into runDir. A missing file yields an empty slice (no
// proposals this run), not an error.
//
// Parameters:
//   - runDir: absolute path to a per-run dreams/<ts>/ directory
//
// Returns:
//   - []Proposal: the proposals the executor emitted, in file order
//   - error: non-nil on a read or JSON decode failure
func ReadProposals(runDir string) ([]Proposal, error) {
	path := filepath.Join(runDir, cfgDream.FileProposals)
	data, readErr := ctxIo.SafeReadUserFile(path)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return []Proposal{}, nil
		}
		return nil, errDream.ReadProposals(path, readErr)
	}
	var proposals []Proposal
	if unmarshalErr := json.Unmarshal(data, &proposals); unmarshalErr != nil {
		return nil, errDream.ReadProposals(path, unmarshalErr)
	}
	return proposals, nil
}

// FindProposal locates the proposal with the given id within runDir.
//
// Parameters:
//   - runDir: absolute path to a per-run dreams/<ts>/ directory
//   - id: the proposal ID to locate
//
// Returns:
//   - Proposal: the matching proposal
//   - error: ProposalNotFound when no proposal carries the id, or a
//     read failure
func FindProposal(runDir, id string) (Proposal, error) {
	proposals, readErr := ReadProposals(runDir)
	if readErr != nil {
		return Proposal{}, readErr
	}
	for _, p := range proposals {
		if p.ID == id {
			return p, nil
		}
	}
	return Proposal{}, errDream.ProposalNotFound(id, runDir)
}

// PendingProposals filters proposals to those not yet recorded in the
// ledger (dedup-against-seen): a proposal whose ID already has a
// disposition is dropped.
//
// Parameters:
//   - proposals: the proposals from a run
//   - ledger: the recorded dispositions (from ReadLedger)
//
// Returns:
//   - []Proposal: proposals with no ledger entry, in input order
func PendingProposals(
	proposals []Proposal, ledger []LedgerEntry,
) []Proposal {
	var pending []Proposal
	for _, p := range proposals {
		if !Seen(ledger, p.ID) {
			pending = append(pending, p)
		}
	}
	return pending
}
