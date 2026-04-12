//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	cfgEntry "github.com/ActiveMemory/ctx/internal/config/entry"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// validateEntry checks a PublishEntry for required fields
// and enforces size limits, including the Meta
// sub-struct.
//
// Parameters:
//   - pe: entry to validate
//
// Returns:
//   - error: non-nil if validation fails
func validateEntry(pe PublishEntry) error {
	if pe.ID == "" {
		return status.Error(
			codes.InvalidArgument, cfgHub.ErrEntryIDRequired,
		)
	}
	if !cfgEntry.AllowedTypes[pe.Type] {
		return status.Errorf(
			codes.InvalidArgument,
			cfgHub.ErrInvalidEntryType, pe.Type,
		)
	}
	if pe.Origin == "" {
		return status.Error(
			codes.InvalidArgument,
			cfgHub.ErrEntryOriginRequired,
		)
	}
	if len(pe.Content) > cfgHub.MaxContentLen {
		return status.Error(
			codes.InvalidArgument,
			cfgHub.ErrEntryContentOversize,
		)
	}
	return validateEntryMeta(pe.Meta)
}

// validateEntryMeta enforces size and character
// restrictions on client-advisory metadata.
//
// Each field is capped at [cfgHub.MaxMetaFieldLen] bytes;
// the sum of all fields is capped at
// [cfgHub.MaxMetaTotalLen]. Fields must be plain
// single-line strings: no newlines, no carriage returns,
// no NUL bytes, no other C0 control characters except
// tab. This prevents log injection (into audits.jsonl),
// markdown injection (into .context/hub/*.md), and
// frontmatter confusion.
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
		{cfgHub.MetaDisplayName, m.DisplayName},
		{cfgHub.MetaHost, m.Host},
		{cfgHub.MetaTool, m.Tool},
		{cfgHub.MetaVia, m.Via},
	}

	total := 0
	for _, f := range fields {
		if len(f.value) > cfgHub.MaxMetaFieldLen {
			return status.Errorf(
				codes.InvalidArgument,
				cfgHub.ErrMetaFieldOversize,
				f.name, cfgHub.MaxMetaFieldLen,
			)
		}
		if metaErr := metaCharCheck(f.name, f.value); metaErr != nil {
			return metaErr
		}
		total += len(f.value)
	}

	if total > cfgHub.MaxMetaTotalLen {
		return status.Errorf(
			codes.InvalidArgument,
			cfgHub.ErrMetaTotalOversize,
			cfgHub.MaxMetaTotalLen,
		)
	}
	return nil
}

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
		if c < cfgHub.MetaControlSpaceLow || c == cfgHub.MetaControlDelete {
			return status.Errorf(
				codes.InvalidArgument,
				cfgHub.ErrMetaControlChar,
				name,
			)
		}
	}
	return nil
}
