//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package relationship

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// renderRow renders one Row as a markdown table row terminated
// by a newline.
//
// Parameters:
//   - row: row content.
//
// Returns:
//   - string: markdown row with trailing newline.
func renderRow(row Row) string {
	var sb strings.Builder
	sb.WriteString(marker.TableRowOpen)
	sb.WriteString(escapeCell(row.From))
	sb.WriteString(marker.TableCellSep)
	sb.WriteString(escapeCell(row.To))
	sb.WriteString(marker.TableCellSep)
	sb.WriteString(escapeCell(row.Kind))
	sb.WriteString(marker.TableCellSep)
	sb.WriteString(escapeCell(row.Summary))
	sb.WriteString(marker.TableRowClose)
	sb.WriteString(token.NewlineLF)
	return sb.String()
}

// escapeCell makes a free-text field safe to embed in a
// markdown table cell.
//
// Parameters:
//   - s: raw cell content.
//
// Returns:
//   - string: escaped cell content.
func escapeCell(s string) string {
	s = strings.ReplaceAll(s, token.NewlineCRLF, token.Space)
	s = strings.ReplaceAll(s, token.NewlineLF, token.Space)
	s = strings.ReplaceAll(
		s, marker.TablePipe, marker.TablePipeEscaped,
	)
	return strings.TrimSpace(s)
}
