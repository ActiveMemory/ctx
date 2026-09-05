//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"bytes"
	"encoding/json"
	"strings"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
)

// parseManifest splits a hooks.json document into its top-level
// keys and the per-event matcher groups under "hooks".
//
// Parameters:
//   - data: JSON document
//
// Returns:
//   - map[string]json.RawMessage: top-level keys (including "hooks")
//   - map[string][]json.RawMessage: event name → matcher groups
//   - error: non-nil when the document or its "hooks" value does
//     not have the expected shape
func parseManifest(data []byte) (
	map[string]json.RawMessage, map[string][]json.RawMessage, error,
) {
	top := map[string]json.RawMessage{}
	if topErr := json.Unmarshal(data, &top); topErr != nil {
		return nil, nil, topErr
	}
	events := map[string][]json.RawMessage{}
	if raw, hasHooks := top[cfgCodex.KeyHooks]; hasHooks {
		if hooksErr := json.Unmarshal(raw, &events); hooksErr != nil {
			return nil, nil, hooksErr
		}
	}
	return top, events, nil
}

// foreignGroups returns the matcher groups that are not
// ctx-managed, in their original order.
//
// Parameters:
//   - groups: matcher groups of one event
//
// Returns:
//   - []json.RawMessage: groups to preserve (nil when none)
func foreignGroups(groups []json.RawMessage) []json.RawMessage {
	var kept []json.RawMessage
	for _, group := range groups {
		if stripped, keep := withoutManagedHandlers(group); keep {
			kept = append(kept, stripped)
		}
	}
	return kept
}

// withoutManagedHandlers removes ctx-managed handlers from a
// matcher group. A pure-ctx group is dropped entirely (the fresh
// embedded groups replace it); a mixed group — the user added
// their own handler next to ctx's — keeps only the user handlers
// so a merge never duplicates the ctx hooks AND never deletes
// user content. Unparseable groups pass through untouched.
//
// Parameters:
//   - group: one matcher group
//
// Returns:
//   - json.RawMessage: the group with ctx handlers removed
//   - bool: false when the group should be dropped
func withoutManagedHandlers(
	group json.RawMessage,
) (json.RawMessage, bool) {
	fields := map[string]json.RawMessage{}
	if groupErr := json.Unmarshal(group, &fields); groupErr != nil {
		return group, true
	}
	var handlers []json.RawMessage
	if handlersErr := json.Unmarshal(
		fields[cfgCodex.KeyHandlers], &handlers,
	); handlersErr != nil || len(handlers) == 0 {
		return group, true
	}

	var foreign []json.RawMessage
	for _, h := range handlers {
		if !managedHandler(h) {
			foreign = append(foreign, h)
		}
	}
	switch {
	case len(foreign) == len(handlers):
		return group, true
	case len(foreign) == 0:
		return nil, false
	}
	encoded, encodeErr := encode(foreign, "")
	if encodeErr != nil {
		return group, true
	}
	fields[cfgCodex.KeyHandlers] = json.RawMessage(
		strings.TrimSpace(string(encoded)),
	)
	rebuilt, rebuildErr := encode(fields, "")
	if rebuildErr != nil {
		return group, true
	}
	return json.RawMessage(strings.TrimSpace(string(rebuilt))), true
}

// managedHandler reports whether one handler's command carries a
// ctx-managed prefix (current or any legacy shape).
//
// Parameters:
//   - handler: one handler object
//
// Returns:
//   - bool: true when the command is ctx-managed
func managedHandler(handler json.RawMessage) bool {
	fields := map[string]json.RawMessage{}
	if hErr := json.Unmarshal(handler, &fields); hErr != nil {
		return false
	}
	var command string
	if cmdErr := json.Unmarshal(
		fields[cfgCodex.KeyCommand], &command,
	); cmdErr != nil {
		return false
	}
	return strings.HasPrefix(command, cfgCodex.HookCommandPrefix) ||
		strings.HasPrefix(
			command, cfgCodex.LegacyHookCommandPrefixGuardless,
		) ||
		strings.HasPrefix(command, cfgCodex.LegacyHookCommandPrefix)
}

// encode marshals a value without HTML escaping, with the given
// indent, and a trailing newline (the json.Encoder contract).
//
// Parameters:
//   - v: value to encode
//   - indent: indent string (empty for compact output)
//
// Returns:
//   - []byte: encoded JSON
//   - error: non-nil on marshal failure
func encode(v any, indent string) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", indent)
	if encodeErr := encoder.Encode(v); encodeErr != nil {
		return nil, encodeErr
	}
	return buf.Bytes(), nil
}

// equivalent reports whether two JSON documents carry the same
// data regardless of key order and whitespace.
//
// Parameters:
//   - a: first document
//   - b: second document
//
// Returns:
//   - bool: true when both decode to the same value
//   - error: non-nil when either document does not parse
func equivalent(a, b []byte) (bool, error) {
	var av, bv any
	if aErr := json.Unmarshal(a, &av); aErr != nil {
		return false, aErr
	}
	if bErr := json.Unmarshal(b, &bv); bErr != nil {
		return false, bErr
	}
	ac, acErr := json.Marshal(av)
	if acErr != nil {
		return false, acErr
	}
	bc, bcErr := json.Marshal(bv)
	if bcErr != nil {
		return false, bcErr
	}
	return bytes.Equal(ac, bc), nil
}
