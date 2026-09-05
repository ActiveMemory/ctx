//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"bytes"
	"encoding/json"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
)

// MergeHooks combines an existing hooks.json with the embedded ctx
// manifest.
//
// Semantics:
//   - empty existing content: the embedded manifest is returned
//     verbatim ([OutcomeCreated]);
//   - per event, foreign matcher groups are kept in their original
//     order, ctx-managed groups (every handler command starts with
//     [cfgCodex.HookAnchor]) are dropped, and the embedded groups
//     are appended; events only present in the existing file are
//     preserved;
//   - other top-level keys are preserved; a missing "description"
//     is taken from the embedded manifest;
//   - when the merged document is semantically identical to the
//     existing one, the existing bytes are returned with
//     [OutcomeSkipped] so callers leave the file alone.
//
// Parameters:
//   - existing: current file content (may be empty)
//   - embedded: the embedded ctx hooks manifest
//
// Returns:
//   - []byte: content to write (2-space indent, trailing newline)
//   - Outcome: Created, Merged, or Skipped
//   - error: non-nil when either input is not a JSON object
func MergeHooks(existing, embedded []byte) ([]byte, Outcome, error) {
	if len(bytes.TrimSpace(existing)) == 0 {
		return embedded, OutcomeCreated, nil
	}

	current, currentEvents, parseErr := parseManifest(existing)
	if parseErr != nil {
		return nil, OutcomeSkipped, errParser.Unmarshal(parseErr)
	}
	shipped, shippedEvents, shippedErr := parseManifest(embedded)
	if shippedErr != nil {
		return nil, OutcomeSkipped, errParser.Unmarshal(shippedErr)
	}

	merged := map[string][]json.RawMessage{}
	for event, groups := range currentEvents {
		merged[event] = foreignGroups(groups)
	}
	for event, groups := range shippedEvents {
		merged[event] = append(merged[event], groups...)
	}
	for event, groups := range merged {
		if len(groups) == 0 {
			delete(merged, event)
		}
	}

	hooksRaw, hooksErr := encode(merged, "")
	if hooksErr != nil {
		return nil, OutcomeSkipped, hooksErr
	}
	current[cfgCodex.KeyHooks] = json.RawMessage(bytes.TrimSpace(hooksRaw))
	if _, hasDesc := current[cfgCodex.KeyDescription]; !hasDesc {
		if description, ok := shipped[cfgCodex.KeyDescription]; ok {
			current[cfgCodex.KeyDescription] = description
		}
	}

	out, encodeErr := encode(current, token.Indent2)
	if encodeErr != nil {
		return nil, OutcomeSkipped, encodeErr
	}

	if same, sameErr := equivalent(existing, out); sameErr == nil && same {
		return existing, OutcomeSkipped, nil
	}
	return out, OutcomeMerged, nil
}
