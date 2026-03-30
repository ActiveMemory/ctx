//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package format

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgFmt "github.com/ActiveMemory/ctx/internal/config/format"
)

// Tokens formats a token count as a human-readable string with SI suffix.
//
// Parameters:
//   - tokens: Token count to format
//
// Returns:
//   - string: Formatted count (e.g., "500", "1.5K", "2.3M")
func Tokens(tokens int) string {
	if tokens < cfgFmt.SIThreshold {
		return fmt.Sprintf(desc.Text(text.DescKeyWriteFormatSIInteger), tokens)
	}
	if tokens < cfgFmt.SIThresholdM {
		return fmt.Sprintf(
			desc.Text(text.DescKeyWriteFormatSIKiloUpper),
			float64(tokens)/cfgFmt.SIThreshold,
		)
	}
	return fmt.Sprintf(
		desc.Text(text.DescKeyWriteFormatSIMegaUpper),
		float64(tokens)/cfgFmt.SIThresholdM,
	)
}
