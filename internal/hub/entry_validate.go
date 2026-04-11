//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// maxContentLen is the maximum entry content size (1MB).
// Entries are text-only (decisions, learnings, conventions).
const maxContentLen = 1 << 20

// Size caps for EntryMeta fields. Deliberately tight
// because Meta is advisory, not content-bearing.
const (
	// maxMetaFieldLen is the per-field size cap.
	maxMetaFieldLen = 256
	// maxMetaTotalLen caps the sum of all Meta field
	// lengths to prevent abuse via many nearly-full
	// fields.
	maxMetaTotalLen = 2048
)

// allowedTypes is the set of valid entry types.
var allowedTypes = map[string]bool{
	"decision":   true,
	"learning":   true,
	"convention": true,
	"task":       true,
}

// validateEntry checks a PublishEntry for required fields
// and enforces size limits, including the Meta
// sub-struct.
func validateEntry(pe PublishEntry) error {
	if pe.ID == "" {
		return status.Error(
			codes.InvalidArgument, "entry ID required",
		)
	}
	if !allowedTypes[pe.Type] {
		return status.Errorf(
			codes.InvalidArgument,
			"invalid entry type %q", pe.Type,
		)
	}
	if pe.Origin == "" {
		return status.Error(
			codes.InvalidArgument,
			"entry origin required",
		)
	}
	if len(pe.Content) > maxContentLen {
		return status.Error(
			codes.InvalidArgument,
			"entry content exceeds 1MB limit",
		)
	}
	return validateEntryMeta(pe.Meta)
}

// validateEntryMeta enforces size and character
// restrictions on client-advisory metadata.
//
// Each field is capped at [maxMetaFieldLen] bytes; the
// sum of all fields is capped at [maxMetaTotalLen].
// Fields must be plain single-line strings: no newlines,
// no carriage returns, no NUL bytes, no other C0 control
// characters except tab. This prevents log injection
// (into audits.jsonl), markdown injection (into
// .context/hub/*.md), and frontmatter confusion.
//
// Parameters:
//   - m: client-sent metadata
//
// Returns:
//   - error: non-nil if any restriction is violated
func validateEntryMeta(m EntryMeta) error {
	fields := []struct {
		name  string
		value string
	}{
		{"display_name", m.DisplayName},
		{"host", m.Host},
		{"tool", m.Tool},
		{"via", m.Via},
	}

	total := 0
	for _, f := range fields {
		if len(f.value) > maxMetaFieldLen {
			return status.Errorf(
				codes.InvalidArgument,
				"meta.%s exceeds %d bytes",
				f.name, maxMetaFieldLen,
			)
		}
		if metaErr := metaCharCheck(f.name, f.value); metaErr != nil {
			return metaErr
		}
		total += len(f.value)
	}

	if total > maxMetaTotalLen {
		return status.Errorf(
			codes.InvalidArgument,
			"meta total exceeds %d bytes",
			maxMetaTotalLen,
		)
	}
	return nil
}

// metaControlSpaceLow is the lowest printable ASCII byte;
// anything below this (except tab) is a control character.
const metaControlSpaceLow = 0x20

// metaControlDelete is the DEL character, the only C1 we
// reject by number rather than by range.
const metaControlDelete = 0x7f

// metaCharCheck rejects control characters other than
// horizontal tab in a Meta field value. Tab is allowed
// because it shows up in shell prompts and tool strings
// and is safe for single-line renderers; newlines,
// carriage returns, NUL, and DEL are rejected to prevent
// log / markdown injection.
//
// Parameters:
//   - name: field name for error reporting
//   - value: candidate string value
//
// Returns:
//   - error: non-nil on a disallowed character
func metaCharCheck(name, value string) error {
	tab := token.Tab[0]
	for i := 0; i < len(value); i++ {
		c := value[i]
		if c == tab {
			continue
		}
		if c < metaControlSpaceLow || c == metaControlDelete {
			return status.Errorf(
				codes.InvalidArgument,
				"meta.%s contains control character",
				name,
			)
		}
	}
	return nil
}
