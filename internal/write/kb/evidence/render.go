//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package evidence

import (
	"fmt"
	"strings"
	"time"

	cfgKbEvidence "github.com/ActiveMemory/ctx/internal/config/kb/evidence"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// renderRow turns a Row into one markdown table line (no
// trailing newline).
//
// Parameters:
//   - r: row to render.
//
// Returns:
//   - string: pipe-delimited table row.
func renderRow(r Row) string {
	tags := strings.Join(r.Tags, token.CommaSpace)
	extracted := r.Extracted.UTC().Format(time.DateOnly)
	return fmt.Sprintf(
		cfgKbEvidence.RowFormat,
		r.ID, escapeCell(r.Claim), r.SourceID, r.Locator,
		r.SHA, r.Confidence, tags, r.Occurred, extracted,
	)
}

// escapeCell escapes pipe characters that would otherwise
// break the markdown table grammar.
//
// Parameters:
//   - s: free text.
//
// Returns:
//   - string: s with `|` replaced by `\|`.
func escapeCell(s string) string {
	return strings.ReplaceAll(
		s, marker.TablePipe, marker.TablePipeEscaped,
	)
}
