//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"strconv"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
)

// parseRef breaks a raw reference string into its components.
//
// Formats:
//
//	"decision:12"  → ("decision", 12, "")
//	"session:abc"  → ("session", 0, "abc")
//	`"Some note"`  → ("note", 0, "Some note")
//	unknown        → ("note", 0, ref)
//
// Parameters:
//   - ref: raw reference string
//
// Returns:
//   - refType: type keyword (decision, learning, convention, task, session, note)
//   - number: numeric value, 0 when not applicable
//   - text: text value, empty when not applicable
func parseRef(ref string) (refType string, number int, text string) {
	// Quoted strings are free-form notes.
	if strings.HasPrefix(ref, token.DoubleQuote) && strings.HasSuffix(ref, token.DoubleQuote) {
		return cfgTrace.RefTypeNote, 0, strings.Trim(ref, token.DoubleQuote)
	}

	parts := strings.SplitN(ref, token.Colon, 2)
	if len(parts) != 2 {
		return cfgTrace.RefTypeNote, 0, ref
	}

	kind := parts[0]
	value := parts[1]

	switch kind {
	case cfgTrace.RefTypeDecision, cfgTrace.RefTypeLearning,
		cfgTrace.RefTypeConvention, cfgTrace.RefTypeTask:
		n, err := strconv.Atoi(value)
		if err != nil {
			return cfgTrace.RefTypeNote, 0, ref
		}
		return kind, n, ""
	case cfgTrace.RefTypeSession:
		return cfgTrace.RefTypeSession, 0, value
	default:
		return cfgTrace.RefTypeNote, 0, ref
	}
}
