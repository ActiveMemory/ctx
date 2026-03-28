//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package frontmatter

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// ResolveHeading returns the first non-empty value among title, slug, baseName.
//
// Parameters:
//   - title: Preferred heading text
//   - slug: Fallback slug from session metadata
//   - baseName: Last-resort filename base
//
// Returns:
//   - string: First non-empty value
func ResolveHeading(title, slug, baseName string) string {
	if title != "" {
		return title
	}
	if slug != "" {
		return slug
	}
	return baseName
}

// WriteFmQuoted writes a YAML frontmatter quoted string field.
//
// Parameters:
//   - sb: String builder to write to
//   - key: Frontmatter key
//   - value: Quoted string value
func WriteFmQuoted(sb *strings.Builder, key, value string) {
	_, writeErr := fmt.Fprintf(sb, tpl.FmQuoted+token.NewlineLF, key, value)
	if writeErr != nil {
		return
	}
}

// WriteFmString writes a YAML frontmatter bare string field.
//
// Parameters:
//   - sb: String builder to write to
//   - key: Frontmatter key
//   - value: Bare string value
func WriteFmString(sb *strings.Builder, key, value string) {
	_, writeErr := fmt.Fprintf(sb, tpl.FmString+token.NewlineLF, key, value)
	if writeErr != nil {
		return
	}
}

// WriteFmInt writes a YAML frontmatter integer field.
//
// Parameters:
//   - sb: String builder to write to
//   - key: Frontmatter key
//   - value: Integer value
func WriteFmInt(sb *strings.Builder, key string, value int) {
	_, writeErr := fmt.Fprintf(sb, tpl.FmInt+token.NewlineLF, key, value)
	if writeErr != nil {
		return
	}
}
