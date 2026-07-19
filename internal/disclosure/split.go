//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"strings"

	cfgDisc "github.com/ActiveMemory/ctx/internal/config/disclosure"
	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
	"github.com/ActiveMemory/ctx/internal/heading"
)

// FlattenPlan resolves a plan's per-theme staged entries into one ordered
// list of entry identities (timestamp+title), in assignment order. An
// assignment that names no entries is a malformed plan (a theme with
// nothing to receive), not a no-op.
//
// Parameters:
//   - p: the digest plan
//
// Returns:
//   - []string: the entry IDs to move, in assignment order
//   - error: ErrEmptyAssignment when an assignment has no entries
func FlattenPlan(p Plan) ([]string, error) {
	var ids []string
	for _, a := range p.Assignments {
		if len(a.Entries) == 0 {
			return nil, errDisc.ErrEmptyAssignment
		}
		for _, e := range a.Entries {
			ids = append(ids, entryID(e))
		}
	}
	return ids, nil
}

// SplitStaging removes the entries named by moveIDs from a raw staging
// zone and returns each moved entry's verbatim span alongside the
// remaining staging, byte-exact.
//
// It cuts on raw header-to-next-header byte spans rather than reusing
// heading.ParseEntryBlocks' BlockContent — that helper trims trailing
// blank lines, so re-joining parsed blocks would silently drop the
// inter-entry whitespace. ParseEntryBlocks is used only to identify the
// entries and their order; the spans are cut from the original bytes so
// Reconstruct-grade fidelity holds (spec Guards §1, Conservation).
//
// Parameters:
//   - staging: the root's raw staging zone (starts at the first "## [")
//   - moveIDs: entry IDs (timestamp+IDSeparator+title) to move out
//
// Returns:
//   - map[string]string: moved id -> its verbatim span
//   - string: the remaining staging (un-moved entries, byte-exact)
//   - error: ErrEntryAssignedTwice (dup id), ErrEntryNotInStaging (unknown)
func SplitStaging(
	staging string, moveIDs []string,
) (map[string]string, string, error) {
	want := make(map[string]bool, len(moveIDs))
	for _, id := range moveIDs {
		if want[id] {
			return nil, "", errDisc.ErrEntryAssignedTwice
		}
		want[id] = true
	}

	blocks := heading.ParseEntryBlocks(staging)
	byteOff := lineByteOffsets(staging)

	moved := make(map[string]string, len(moveIDs))
	seen := make(map[string]bool, len(moveIDs))
	// Removal spans, collected in increasing byte order (blocks are sorted).
	type span struct{ start, end int }
	var removals []span

	for i, b := range blocks {
		id := b.Entry.Timestamp + cfgDisc.IDSeparator + b.Entry.Title
		if !want[id] || seen[id] {
			continue
		}
		endLine := len(byteOff) - 1
		if i+1 < len(blocks) {
			endLine = blocks[i+1].StartIndex
		}
		start := clampOffset(byteOff[b.StartIndex], len(staging))
		end := clampOffset(byteOff[endLine], len(staging))
		moved[id] = staging[start:end]
		removals = append(removals, span{start, end})
		seen[id] = true
	}

	for _, id := range moveIDs {
		if !seen[id] {
			return nil, "", errDisc.ErrEntryNotInStaging
		}
	}

	var sb strings.Builder
	prev := 0
	for _, r := range removals {
		sb.WriteString(staging[prev:r.start])
		prev = r.end
	}
	sb.WriteString(staging[prev:])
	return moved, sb.String(), nil
}
