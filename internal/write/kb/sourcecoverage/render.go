//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sourcecoverage

import (
	"strings"
	"time"

	cfgKbSC "github.com/ActiveMemory/ctx/internal/config/kb/sourcecoverage"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// render composes the ledger file: title + lead paragraphs +
// table header + one row per source.
//
// Parameters:
//   - rows: rows to render.
//
// Returns:
//   - string: full file content.
func render(rows []Row) string {
	var sb strings.Builder
	sb.WriteString(cfgKbSC.TitleHeading)
	sb.WriteString(token.DoubleNewline)
	sb.WriteString(cfgKbSC.LeadParagraph1)
	sb.WriteString(token.NewlineLF)
	sb.WriteString(cfgKbSC.LeadParagraph2)
	sb.WriteString(token.NewlineLF)
	sb.WriteString(cfgKbSC.LeadParagraph3)
	sb.WriteString(token.DoubleNewline)
	sb.WriteString(cfgKbSC.TableHeader)
	sb.WriteString(token.NewlineLF)
	for _, r := range rows {
		sb.WriteString(marker.TableRowOpen)
		sb.WriteString(r.Source)
		sb.WriteString(marker.TableCellSep)
		sb.WriteString(r.Topic)
		sb.WriteString(marker.TableCellSep)
		sb.WriteString(r.State)
		sb.WriteString(marker.TableCellSep)
		sb.WriteString(r.EVCoverage)
		sb.WriteString(marker.TableCellSep)
		sb.WriteString(r.Residue)
		sb.WriteString(marker.TableCellSep)
		sb.WriteString(r.NextAction)
		sb.WriteString(marker.TableCellSep)
		sb.WriteString(r.Updated.UTC().Format(time.DateOnly))
		sb.WriteString(marker.TableRowClose)
		sb.WriteString(token.NewlineLF)
	}
	return sb.String()
}
