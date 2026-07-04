//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package statusline

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	core "github.com/ActiveMemory/ctx/internal/cli/system/core/statusline"
	cfgStatusline "github.com/ActiveMemory/ctx/internal/config/statusline"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeStatusline "github.com/ActiveMemory/ctx/internal/write/statusline"
)

// Run executes the statusline rendering logic.
//
// Reads the status line JSON payload from stdin and prints one
// sanitized line assembled from location, model, context-usage, and
// cost segments. Missing fields drop their segment; malformed input
// degrades to whatever remains renderable. When statusline.enabled
// is false in .ctxrc, prints an empty line so the displayed status
// line goes blank immediately.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input carrying the status line JSON payload
//
// Returns:
//   - error: Always nil (a non-zero exit blanks the status line)
func Run(cmd *cobra.Command, stdin *os.File) error {
	if !rc.StatuslineEnabled() {
		// statusline.enabled: false takes effect immediately by
		// rendering nothing; the settings.local.json entry itself is
		// restored/removed the next time the init merge runs.
		writeStatusline.Blank(cmd)
		return nil
	}
	var p core.Payload
	if raw, readErr := io.ReadAll(
		io.LimitReader(stdin, cfgStatusline.MaxPayloadBytes),
	); readErr == nil {
		// Unmarshal errors are deliberately ignored: a partial or
		// malformed payload still renders a degraded line below.
		_ = json.Unmarshal(raw, &p)
	}

	segments := make([]string, 0, cfgStatusline.SegmentCapacity)
	if loc := core.LocationSegment(&p); loc != "" {
		segments = append(segments, loc)
	}
	if model := core.Sanitize(p.Model.DisplayName); model != "" {
		segments = append(segments, model)
	}
	if pct := p.ContextWindow.UsedPercentage; pct != nil &&
		*pct >= 0 && *pct <= cfgStatusline.PercentMax {
		segments = append(segments,
			fmt.Sprintf(cfgStatusline.ContextFormat, *pct))
	}
	if cost := p.Cost.TotalCostUSD; cost != nil && *cost >= 0 &&
		rc.StatuslineShowCost() {
		segments = append(segments,
			fmt.Sprintf(cfgStatusline.CostFormat, *cost))
	}

	line := strings.Join(segments, cfgStatusline.SegmentSeparator)
	if len(line) > cfgStatusline.MaxLineLen {
		line = line[:cfgStatusline.MaxLineLen]
	}
	writeStatusline.Line(cmd, line)
	return nil
}
