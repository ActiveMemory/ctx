//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
)

// Resolve looks up a raw reference and returns its full details.
//
// Parameters:
//   - ref: raw reference string (e.g. "decision:12", "task:8", `"Some note"`)
//   - contextDir: absolute path to the .context/ directory
//
// Returns:
//   - ResolvedRef: resolved reference with title, detail, and found status
func Resolve(ref, contextDir string) ResolvedRef {
	refType, number, text := parseRef(ref)

	resolved := ResolvedRef{
		Raw:    ref,
		Type:   refType,
		Number: number,
	}

	switch refType {
	case cfgTrace.RefTypeDecision:
		return resolveEntry(resolved, contextDir, cfgCtx.Decision, number)
	case cfgTrace.RefTypeLearning:
		return resolveEntry(resolved, contextDir, cfgCtx.Learning, number)
	case cfgTrace.RefTypeConvention:
		return resolveEntry(resolved, contextDir, cfgCtx.Convention, number)
	case cfgTrace.RefTypeTask:
		return resolveTask(resolved, contextDir, number)
	case cfgTrace.RefTypeSession:
		resolved.Title = text
		resolved.Found = true
		return resolved
	default: // cfgTrace.RefTypeNote
		resolved.Title = text
		resolved.Found = true
		return resolved
	}
}
