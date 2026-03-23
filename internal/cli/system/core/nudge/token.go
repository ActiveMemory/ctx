//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package nudge

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/stats"
)

// TokenUsageLine formats a context window usage line for display.
// Shows an icon (normal or warning), token count, percentage, and window size.
//
// Parameters:
//   - tokens: number of tokens used
//   - pct: percentage of the context window used
//   - windowSize: total context window size
//
// Returns:
//   - string: formatted usage line (e.g., "⏱ Context window: ~12k tokens (~60% of 200k)")
func TokenUsageLine(tokens, pct, windowSize int) string {
	icon := desc.Text(text.DescKeyCheckContextSizeTokenNormal)
	suffix := ""
	if pct >= stats.ContextWindowThresholdPct {
		icon = desc.Text(text.DescKeyCheckContextSizeTokenLow)
		suffix = desc.Text(text.DescKeyCheckContextSizeRunningLowSuffix)
	}
	return fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeTokenUsage),
		icon, coreSession.FormatTokenCount(tokens),
		pct, coreSession.FormatWindowSize(windowSize), suffix,
	)
}
