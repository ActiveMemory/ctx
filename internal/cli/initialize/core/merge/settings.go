//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	"bytes"
	"encoding/json"

	cfgClaude "github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/err/config"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// readSettingsRaw reads settings.local.json into a key-level raw map.
//
// The raw-map representation exists so read-modify-write cycles only
// touch the sections they own: keys ctx does not model (a user's env,
// model, or other settings) survive the rewrite byte-for-byte instead
// of being silently dropped by a typed round-trip.
//
// A missing or unreadable file is reported as non-existent with an
// empty map, matching the merge semantics of a fresh project.
//
// Returns:
//   - map[string]json.RawMessage: top-level settings keys
//   - bool: true when the settings file exists and was read
//   - error: non-nil when the file exists but does not parse
func readSettingsRaw() (map[string]json.RawMessage, bool, error) {
	raw := map[string]json.RawMessage{}
	content, readErr := ctxIo.SafeReadUserFile(cfgClaude.Settings)
	if readErr != nil {
		return raw, false, nil
	}
	if unmarshalErr := json.Unmarshal(content, &raw); unmarshalErr != nil {
		return nil, true, errParser.ParseFile(cfgClaude.Settings, unmarshalErr)
	}
	return raw, true, nil
}

// writeSettingsRaw writes the raw settings map back to
// settings.local.json, creating .claude/ when needed.
//
// Parameters:
//   - raw: top-level settings keys to persist
//
// Returns:
//   - error: non-nil if encoding or file operations fail
func writeSettingsRaw(raw map[string]json.RawMessage) error {
	if mkdirErr := ctxIo.SafeMkdirAll(dir.Claude, fs.PermExec); mkdirErr != nil {
		return errFs.Mkdir(dir.Claude, mkdirErr)
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", token.Indent2)
	if encodeErr := encoder.Encode(raw); encodeErr != nil {
		return config.MarshalSettings(encodeErr)
	}
	if writeErr := ctxIo.SafeWriteFile(
		cfgClaude.Settings, buf.Bytes(), fs.PermFile,
	); writeErr != nil {
		return errFs.FileWrite(cfgClaude.Settings, writeErr)
	}
	return nil
}

// marshalSettingsSection marshals one settings section without HTML
// escaping, for insertion into the raw settings map.
//
// Parameters:
//   - v: section value to marshal
//
// Returns:
//   - json.RawMessage: marshaled section
//   - error: non-nil on marshal failure
func marshalSettingsSection(v any) (json.RawMessage, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if encodeErr := encoder.Encode(v); encodeErr != nil {
		return nil, encodeErr
	}
	return json.RawMessage(bytes.TrimSpace(buf.Bytes())), nil
}
