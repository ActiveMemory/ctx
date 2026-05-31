//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package jsonpayload

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/spf13/cobra"

	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errAdd "github.com/ActiveMemory/ctx/internal/err/add"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Load reads and strictly decodes a JSON payload file.
//
// Decoding rejects unknown keys so that a misspelled field surfaces as
// an error instead of being silently dropped.
//
// Parameters:
//   - path: Filesystem path to the JSON payload
//
// Returns:
//   - Payload: Decoded payload
//   - error: Non-nil if the file cannot be read or the JSON is invalid
func Load(path string) (Payload, error) {
	var p Payload

	data, readErr := ctxIo.SafeReadUserFile(path)
	if readErr != nil {
		return p, errFs.FileRead(path, readErr)
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if decErr := dec.Decode(&p); decErr != nil {
		return p, errAdd.JSONParse(path, decErr)
	}
	return p, nil
}

// Content returns the entry content derived from the payload.
//
// The content is the trimmed Title; for entries that carry a Body
// (tasks), a non-empty Body is space-joined onto the Title since the
// target file stores each entry on a single line. Returns an empty
// string when neither field is set, so callers can fall through to
// other content sources.
//
// Returns:
//   - string: The space-joined content, or "" when empty
func (p Payload) Content() string {
	var parts []string
	if title := strings.TrimSpace(p.Title); title != "" {
		parts = append(parts, title)
	}
	if body := strings.TrimSpace(p.Body); body != "" {
		parts = append(parts, body)
	}
	return strings.Join(parts, token.Space)
}

// OverlayFlags loads the --json-file payload (if any) and writes its
// non-empty typed fields onto the command's flags, so downstream
// validation and the bound flag variables see the effective values.
//
// JSON values supersede individually-supplied flags, per spec. Only
// flags that exist on the command and whose payload value is non-empty
// are set; absent payload fields leave any CLI-supplied flag untouched.
// A no-op (returns nil) when --json-file is unset.
//
// Wiring is intentionally explicit: each add noun calls this at the top
// of its PreRunE so the overlay is visible at the call site.
//
// Parameters:
//   - cmd: The cobra command whose flags receive the overlay
//
// Returns:
//   - error: Non-nil if the payload cannot be loaded or a flag set fails
func OverlayFlags(cmd *cobra.Command) error {
	flags := cmd.Flags()

	path, getErr := flags.GetString(cFlag.JSONFile)
	if getErr != nil {
		return getErr
	}
	if path == "" {
		return nil
	}

	p, loadErr := Load(path)
	if loadErr != nil {
		return loadErr
	}

	overlay := []struct {
		name  string
		value string
	}{
		{cFlag.Context, p.Context},
		{cFlag.Rationale, p.Rationale},
		{cFlag.Consequence, p.Consequence},
		{cFlag.Lesson, p.Lesson},
		{cFlag.Application, p.Application},
		{cFlag.Priority, p.Priority},
		{cFlag.Section, p.Section},
		{cFlag.SessionID, p.Provenance.SessionID},
		{cFlag.Branch, p.Provenance.Branch},
		{cFlag.Commit, p.Provenance.Commit},
	}
	for _, f := range overlay {
		if f.value == "" || flags.Lookup(f.name) == nil {
			continue
		}
		if setErr := flags.Set(f.name, f.value); setErr != nil {
			return setErr
		}
	}
	return nil
}
