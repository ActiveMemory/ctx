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

// parse turns the raw ledger contents into rows. Lines that do
// not look like table rows are silently skipped so freeform
// prose at the top of the file is tolerated.
//
// Parameters:
//   - raw: file contents.
//
// Returns:
//   - []Row: parsed rows.
func parse(raw string) []Row {
	var out []Row
	for _, line := range strings.Split(raw, token.NewlineLF) {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, marker.TablePipe) ||
			strings.HasPrefix(line, cfgKbSC.DelimRowPrefix) {
			continue
		}
		cols := splitRow(line)
		if len(cols) != cfgKbSC.ExpectedCellCount {
			continue
		}
		if cols[cfgKbSC.ColSource] == cfgKbSC.HeaderCellSource {
			continue
		}
		updated, _ := time.Parse(
			time.DateOnly, cols[cfgKbSC.ColUpdated],
		)
		out = append(out, Row{
			Source:     cols[cfgKbSC.ColSource],
			Topic:      cols[cfgKbSC.ColTopic],
			State:      cols[cfgKbSC.ColState],
			EVCoverage: cols[cfgKbSC.ColEVCoverage],
			Residue:    cols[cfgKbSC.ColResidue],
			NextAction: cols[cfgKbSC.ColNextAction],
			Updated:    updated,
		})
	}
	return out
}

// splitRow splits a markdown table row into trimmed cell
// contents (without the leading / trailing pipe).
//
// Parameters:
//   - line: one markdown table row (must start with "|").
//
// Returns:
//   - []string: trimmed cell contents.
func splitRow(line string) []string {
	parts := strings.Split(line, marker.TablePipe)
	if len(parts) < 2 {
		return nil
	}
	parts = parts[1 : len(parts)-1]
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}
